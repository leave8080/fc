
const userFunction = const sharp = require('sharp');

exports.resizeImage = async (event) => {
    // 从事件中获取上传图片的相关信息，比如存储路径等
    const imagePath = event.imagePath; 
    // 调用sharp库对图片进行尺寸缩小处理
    await sharp(imagePath)
      .resize(300, 300)
      .toFile('resized_' + imagePath);
    return {
        message: '图片处理完成'
    };
};;

// 从环境变量读取输入
const eventStr = process.env.FUNCTION_EVENT || '{}';
const contextStr = process.env.FUNCTION_CONTEXT || '{}';

const event = JSON.parse(eventStr);
const context = JSON.parse(contextStr);

// 执行用户函数
async function execute() {
    try {
        const result = await userFunction(event, context);
        console.log(JSON.stringify(result));
    } catch (error) {
        console.error(JSON.stringify({error: error.message}));
        process.exit(1);
    }
}

execute();
