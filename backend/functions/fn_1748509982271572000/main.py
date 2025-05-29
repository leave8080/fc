
import json
import os
import sys

def convert_to_uppercase(event, context):\n    \"\"\"\n    将输入文本转换为大写\n    \"\"\"\n    # 检查输入事件是否包含text字段\n    if not isinstance(event, dict):\n        return {\n            \"error\": \"输入必须是字典格式\",\n            \"success\": False\n        }\n    \n    text = event.get('text', '')\n    \n    if not text:\n        return {\n            \"error\": \"缺少text字段或text为空\",\n            \"success\": False\n        }\n    \n    # 转换为大写\n    uppercase_text = text.upper()\n    \n    # 统计信息\n    original_length = len(text)\n    char_count = len([c for c in text if c.isalpha()])\n    digit_count = len([c for c in text if c.isdigit()])\n    space_count = text.count(' ')\n    \n    return {\n        \"success\": True,\n        \"original_text\": text,\n        \"uppercase_text\": uppercase_text,\n        \"statistics\": {\n            \"total_length\": original_length,\n            \"letter_count\": char_count,\n            \"digit_count\": digit_count,\n            \"space_count\": space_count\n        },\n        \"processed_by\": f\"Python函数 (用户: {context.get('user', 'unknown')})\",\n        \"operation\": \"大写转换\"\n    }

def main():
    try:
        # 从环境变量读取输入
        event_str = os.environ.get('FUNCTION_EVENT', '{}')
        context_str = os.environ.get('FUNCTION_CONTEXT', '{}')
        
        event = json.loads(event_str)
        context = json.loads(context_str)
        
        # 执行用户函数
        result = convert_to_uppercase(event, context)
        
        print(json.dumps(result))
    except Exception as e:
        print(json.dumps({"error": str(e)}), file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
