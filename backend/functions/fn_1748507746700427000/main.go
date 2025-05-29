
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

import (
    "fmt"
    "math"
    "sort"
)

func analyzeData(event interface{}, context map[string]string) interface{} {
    // 解析输入事件
    eventMap, ok := event.(map[string]interface{})
    if !ok {
        return map[string]interface{}{"error": "Invalid event format"}
    }

    // 获取数据数组
    dataInterface, exists := eventMap["numbers"]
    if !exists {
        return map[string]interface{}{"error": "Missing numbers field"}
    }

    // 转换为float64切片
    dataSlice, ok := dataInterface.([]interface{})
    if !ok {
        return map[string]interface{}{"error": "Numbers must be an array"}
    }

    var numbers []float64
    for _, v := range dataSlice {
        if num, ok := v.(float64); ok {
            numbers = append(numbers, num)
        } else {
            return map[string]interface{}{"error": "All elements must be numbers"}
        }
    }

    if len(numbers) == 0 {
        return map[string]interface{}{"error": "Empty numbers array"}
    }

    // 计算统计信息
    result := calculateStats(numbers)
    result["operation"] = "data_analysis"
    result["processed_by"] = fmt.Sprintf("Go Cloud Function (User: %s)", context["user"])
    
    return result
}

func calculateStats(numbers []float64) map[string]interface{} {
    n := len(numbers)
    
    // 求和
    sum := 0.0
    for _, num := range numbers {
        sum += num
    }
    
    // 平均值
    mean := sum / float64(n)
    
    // 排序用于中位数和四分位数
    sorted := make([]float64, n)
    copy(sorted, numbers)
    sort.Float64s(sorted)
    
    // 中位数
    var median float64
    if n%2 == 0 {
        median = (sorted[n/2-1] + sorted[n/2]) / 2
    } else {
        median = sorted[n/2]
    }
    
    // 最大值和最小值
    min := sorted[0]
    max := sorted[n-1]
    
    // 方差和标准差
    variance := 0.0
    for _, num := range numbers {
        variance += math.Pow(num-mean, 2)
    }
    variance /= float64(n)
    stdDev := math.Sqrt(variance)
    
    // 四分位数
    q1 := sorted[n/4]
    q3 := sorted[3*n/4]
    
    return map[string]interface{}{
        "count": n,
        "sum": sum,
        "mean": mean,
        "median": median,
        "min": min,
        "max": max,
        "range": max - min,
        "variance": variance,
        "std_deviation": stdDev,
        "q1": q1,
        "q3": q3,
        "iqr": q3 - q1,
        "sorted_data": sorted,
    }
}

func main() {
	// 从环境变量读取输入
	eventStr := os.Getenv("FUNCTION_EVENT")
	contextStr := os.Getenv("FUNCTION_CONTEXT")
	
	var event interface{}
	var context map[string]string
	
	if eventStr != "" {
		json.Unmarshal([]byte(eventStr), &event)
	}
	if contextStr != "" {
		json.Unmarshal([]byte(contextStr), &context)
	}
	
	// 调用用户函数
	result := analyzeData(event, context)
	
	// 输出结果
	resultBytes, _ := json.Marshal(result)
	fmt.Print(string(resultBytes))
}
