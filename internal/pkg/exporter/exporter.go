package exporter

import (
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Exporter 结构体用于处理CSV导出逻辑
type Exporter struct {
	w            *csv.Writer  // CSV写入器
	ctx          *gin.Context // Gin上下文，主要用于写入HTTP响应
	headers      []string     // CSV表头
	fieldIndices []int        // 结构体中与表头对应的字段索引，按表头顺序排列
	wroteHeader  bool         // 标记表头是否已写入
	err          error        // 记录导出过程中发生的第一个错误
}

// New 创建一个新的Exporter实例
// 它会设置必要的HTTP头以便浏览器下载CSV文件
// filename 参数是下载文件的建议名称
func New(ctx *gin.Context, filename string) *Exporter {
	if filename == "" {
		filename = "export.csv" // 默认文件名
	}
	// 设置HTTP头
	ctx.Header("Content-Type", "text/csv; charset=utf-8") // 指定内容类型和字符集
	// 添加UTF-8 BOM标记，帮助Excel等软件正确识别编码，避免中文乱码
	_, _ = ctx.Writer.Write([]byte{0xEF, 0xBB, 0xBF})                   // 忽略BOM写入的错误，尽力而为
	ctx.Header("Content-Disposition", "attachment; filename="+filename) // 提示浏览器下载

	return &Exporter{
		w:   csv.NewWriter(ctx.Writer), // gin.Context.Writer 实现了 io.Writer
		ctx: ctx,
	}
}

// SetHeader 显式设置CSV的表头行
// 如果不调用此方法，表头将在第一次调用WriteData时根据结构体字段自动生成
// 返回Exporter指针，以便进行链式调用
func (e *Exporter) SetHeader(headers []string) *Exporter {
	if e.err != nil {
		return e // 如果之前已发生错误，则直接返回
	}
	if e.wroteHeader {
		e.err = errors.New("CSV表头已写入，无法再次设置")
		return e
	}
	if len(headers) == 0 {
		e.err = errors.New("显式设置表头时，表头切片不能为空")
		return e
	}

	// 直接写入自定义表头，此时不进行字段映射检查，将在WriteData中进行
	// 如果希望在SetHeader时就验证表头与后续数据类型的匹配性，会更复杂，
	// 因为此时可能还不知道具体的数据类型。
	e.headers = headers // 存储用户设置的表头
	// 注意：此时不实际写入CSV，也不设置 e.wroteHeader = true
	// 表头的实际写入和 fieldIndices 的计算推迟到 WriteData 中，
	// 因为需要结合实际数据类型来映射表头到字段。
	return e
}

// buildStructFieldInfo 分析结构体类型，提取用于CSV导出的字段信息
// 它返回:
// - nameToIdxMap: 一个映射，键是字段名或`csv`标签名，值是字段在结构体中的索引
// - defaultHeaders: 根据字段名或`csv`标签自动生成的默认表头切片
// - defaultIndices: 与defaultHeaders对应的字段索引切片
func buildStructFieldInfo(structType reflect.Type) (nameToIdxMap map[string]int, defaultHeaders []string, defaultIndices []int) {
	nameToIdxMap = make(map[string]int)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// PkgPath 为空字符串表示字段是导出的（首字母大写）
		if field.PkgPath != "" {
			continue // 跳过未导出的字段
		}

		tag := field.Tag.Get("csv")
		if tag == "-" {
			continue // 跳过标记为"-"的字段
		}

		fieldName := field.Name
		headerName := tag
		if headerName == "" {
			headerName = fieldName // 如果没有csv标签，默认使用字段名作为表头
		}

		// 确保一个字段索引只被添加一次到 defaultHeaders 和 defaultIndices
		// 并且优先使用csv标签名作为默认表头
		if _, exists := nameToIdxMap[headerName]; !exists { // 如果headerName还没被用作key（主要防止字段名和tag名一样时重复）
			if _, fieldNameMapped := nameToIdxMap[fieldName]; !fieldNameMapped || headerName == fieldName {
				defaultHeaders = append(defaultHeaders, headerName)
				defaultIndices = append(defaultIndices, i)
			}
		}

		// 建立映射关系，优先使用标签名，但也包含字段名（如果不同）
		nameToIdxMap[headerName] = i
		if headerName != fieldName { // 如果标签名和字段名不同，也添加字段名的映射
			nameToIdxMap[fieldName] = i
		}
	}
	return
}

// WriteData 将结构体切片数据写入CSV
// data 参数必须是一个结构体切片 (例如: []MyStruct) 或指向结构体的指针切片 (例如: []*MyStruct)
// 如果之前没有通过SetHeader()设置表头，则会自动根据结构体字段名或`csv`标签生成表头
// 字段可以通过结构体标签 `csv:"-"` 来跳过导出
// 可以通过 `csv:"自定义列名"` 来指定自定义的列名
// 返回Exporter指针，以便进行链式调用
func (e *Exporter) WriteData(data interface{}) *Exporter {
	if e.err != nil {
		return e // 如果之前已发生错误，则直接返回
	}

	sliceVal := reflect.ValueOf(data)
	if sliceVal.Kind() != reflect.Slice {
		e.err = errors.New("WriteData 的参数必须是一个切片")
		return e
	}

	if sliceVal.Len() == 0 {
		// 如果数据为空，但用户通过SetHeader设置了表头且尚未写入，则写入表头
		if len(e.headers) > 0 && !e.wroteHeader {
			if err := e.w.Write(e.headers); err != nil {
				e.err = fmt.Errorf("写入预设表头失败（空数据情况）: %w", err)
			} else {
				e.wroteHeader = true
				// 注意：此时无法确定fieldIndices，因为没有数据元素类型信息
				// 如果后续再调用WriteData传入非空数据，届时会根据e.headers和新数据类型计算fieldIndices
			}
		}
		return e // 没有数据行可写
	}

	elemType := sliceVal.Type().Elem() // 获取切片元素的类型
	isPointerToStruct := false
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem() // 如果是结构体指针，获取其指向的结构体类型
		isPointerToStruct = true
	}

	if elemType.Kind() != reflect.Struct {
		e.err = errors.New("WriteData 的参数必须是一个结构体切片或结构体指针切片")
		return e
	}

	// 初始化或验证表头和字段索引
	if !e.wroteHeader {
		nameToIdxMap, derivedHeaders, derivedIndices := buildStructFieldInfo(elemType)

		if len(e.headers) == 0 { // 用户未调用SetHeader，或调用时传入空切片
			if len(derivedHeaders) == 0 {
				e.err = errors.New("结构体中没有可导出的字段（或所有字段都被标记为`csv:\"-\"`）")
				return e
			}
			e.headers = derivedHeaders
			e.fieldIndices = derivedIndices
		} else { // 用户通过SetHeader设置了表头，现在根据这些表头和结构体类型来确定字段索引
			var finalIndices []int
			for _, userHeader := range e.headers {
				idx, found := nameToIdxMap[userHeader]
				if !found {
					// 用户指定的表头在结构体中找不到对应的字段（或csv标签）
					// 选项1: 报错 (当前实现)
					// 选项2: 为该列写入空值 (需要修改逻辑，例如 finalIndices = append(finalIndices, -1))
					// 选项3: 跳过该表头 (从e.headers中移除，并调整finalIndices)
					e.err = fmt.Errorf("通过SetHeader指定的表头 '%s' 在结构体中找不到对应的可导出字段或csv标签", userHeader)
					return e
				}
				finalIndices = append(finalIndices, idx)
			}
			e.fieldIndices = finalIndices
		}

		if err := e.w.Write(e.headers); err != nil {
			e.err = fmt.Errorf("写入表头行失败: %w", err)
			return e
		}
		e.wroteHeader = true
	} else if len(e.fieldIndices) == 0 && len(e.headers) > 0 {
		// 表头已写入（可能因为前一个WriteData调用是空数据但有SetHeader），但fieldIndices未设置
		// (例如，第一次调用WriteData([])，第二次调用WriteData(realData))
		// 此时需要根据已写入的e.headers和当前数据类型elemType来计算fieldIndices
		nameToIdxMap, _, _ := buildStructFieldInfo(elemType)
		var finalIndices []int
		for _, writtenHeader := range e.headers {
			idx, found := nameToIdxMap[writtenHeader]
			if !found {
				e.err = fmt.Errorf("先前写入的表头 '%s' 与当前数据批次的结构体字段不匹配", writtenHeader)
				return e
			}
			finalIndices = append(finalIndices, idx)
		}
		e.fieldIndices = finalIndices
		if len(e.fieldIndices) == 0 && len(e.headers) > 0 { // 双重检查
			e.err = errors.New("无法将已写入的表头映射到结构体字段以供数据写入")
			return e
		}
	}

	// 写入数据行
	for i := 0; i < sliceVal.Len(); i++ {
		structVal := sliceVal.Index(i) // 获取切片中的第i个元素 (Value类型)

		if isPointerToStruct {
			if structVal.IsNil() {
				// 如果是指针切片且当前元素为nil，则写入一行空字符串
				emptyRow := make([]string, len(e.fieldIndices))
				if err := e.w.Write(emptyRow); err != nil {
					e.err = fmt.Errorf("为nil结构体指针写入空行失败 (索引 %d): %w", i, err)
					return e
				}
				continue // 处理下一个元素
			}
			structVal = structVal.Elem() // 解引用，获取实际的结构体Value
		}

		var record []string
		for _, fieldIndex := range e.fieldIndices { // 严格按照表头顺序（fieldIndices的顺序）提取字段值
			fieldVal := structVal.Field(fieldIndex)
			record = append(record, e.formatField(fieldVal))
		}

		if err := e.w.Write(record); err != nil {
			e.err = fmt.Errorf("写入数据行 %d 失败: %w", i, err)
			return e // 发生写入错误时停止
		}
	}
	return e
}

// formatField 将reflect.Value转换为其CSV字符串表示形式
func (e *Exporter) formatField(fieldVal reflect.Value) string {
	// 处理指针类型：如果是指针，获取它指向的元素
	if fieldVal.Kind() == reflect.Ptr {
		if fieldVal.IsNil() {
			return "" // nil指针对应空字符串，或可自定义为 "<nil>" 等
		}
		fieldVal = fieldVal.Elem() // 解引用指针
	}

	// 优先检查是否实现了fmt.Stringer接口 (例如自定义类型)
	// 但要小心，像time.Time也实现了Stringer，我们可能想用特定格式处理它
	if fieldVal.IsValid() && fieldVal.CanInterface() {
		if _, isTime := fieldVal.Interface().(time.Time); !isTime { // 排除time.Time，它有特殊处理
			if stringer, ok := fieldVal.Interface().(fmt.Stringer); ok {
				// 确保不是在nil指针接收器上调用String() (尽管上面已经解引用)
				// 这一层检查主要是针对非指针类型或已解引用的指针实现了Stringer的情况
				return stringer.String()
			}
		}
	}

	switch fieldVal.Kind() {
	case reflect.Invalid: // 无效的Value，理论上不应发生在此处
		return ""
	case reflect.String:
		return fieldVal.String()
	case reflect.Bool:
		return strconv.FormatBool(fieldVal.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fieldVal.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(fieldVal.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(fieldVal.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(fieldVal.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", fieldVal.Complex()) // 复数的默认格式
	case reflect.Struct:
		// 特殊处理time.Time结构体
		if t, ok := fieldVal.Interface().(time.Time); ok {
			return t.Format(time.RFC3339) // 使用ISO 8601标准格式 (例如: 2006-01-02T15:04:05Z07:00)
		}
		// 对于其他结构体，CSV通常不直接支持嵌套。可以考虑序列化为JSON字符串，
		// 或者返回一个占位符，或者根据业务需求进行扁平化处理。
		// return fmt.Sprintf("%+v", fieldVal.Interface()) // 例如，使用%+v打印字段
		return "[内嵌结构体]" // 当前返回占位符
	case reflect.Slice, reflect.Array:
		// 切片或数组，CSV单元格通常不直接表示。可以考虑用特定分隔符连接元素，
		// 或序列化为JSON字符串，或返回占位符。
		return "[切片/数组]" // 当前返回占位符
	case reflect.Map:
		// Map类型，同上，CSV单元格不直接表示。
		return "[Map]" // 当前返回占位符
	default:
		// 对于其他未明确处理的类型，尝试使用fmt.Sprintf转换
		if fieldVal.IsValid() && fieldVal.CanInterface() {
			return fmt.Sprintf("%v", fieldVal.Interface())
		}
		return "" // 默认返回空字符串
	}
}

// Flush 将所有缓冲数据写入到底层的io.Writer (即gin.Context.Writer)
// 在所有数据都通过WriteData写入后，务必调用此方法
// 返回在写入或刷新过程中发生的任何错误
func (e *Exporter) Flush() error {
	if e.err != nil {
		// 如果在Flush之前已经发生了错误，尝试刷新缓冲区，但优先返回之前记录的错误
		e.w.Flush() // 尝试刷新，以便任何已成功写入部分能送出
		return e.err
	}
	e.w.Flush()        // 执行刷新操作
	return e.w.Error() // 返回csv.Writer在刷新时可能遇到的错误 (例如客户端断开连接)
}

// Error 返回在导出过程中遇到的第一个错误
// 可以在一系列链式操作后调用此方法检查最终状态
func (e *Exporter) Error() error {
	return e.err
}
