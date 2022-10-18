const template = require("art-template");
const fs = require("fs");
const readline = require("readline");
const path = require("path");

const rd = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

//默认的文件打开注册表地址
const DIR_PROJECT = "HKEY_CLASSES_ROOT\\Directory\\shell\\ToolBoxProject";
//默认的文件背景注册表地址
const DIR_BACKGROUND = "HKEY_CLASSES_ROOT\\Directory\\Background\\shell\\ToolBoxBackground";
//icon文件夹
const ICON_DIR = "ico";

const TOOL_BOX_ICO = "icon.ico";


rd.question("input your absolute scripts path:", answer => {
    try {
        const scriptsDir = answer;
        const icoDir = path.join(answer, ICON_DIR);

        let scripts = fs.readdirSync(scriptsDir);//读取脚本列表

        let subCommands = [];//子命令列表
        let subCommandList = [];//子命令对象列表

        for (const script of scripts) {
            let scriptPath = path.join(scriptsDir, script);
            if (fs.lstatSync(scriptPath).isFile() && script.includes(".cmd")) {
                let scriptName = script.replace(".cmd", "");//脚本名称
                let name = `${scriptName.substring(0, 1).toUpperCase()}${scriptName.substring(1)}`;//注册表名称
                let regKey = `JetBrain${name}`;//注册表key
                let display = `"Open ${name} Here"`;//展示内容
                let iconPath = `"${path.join(icoDir, `${scriptName}.ico`)}"`.replace(/\\/g, "\\\\");//icon地址
                let command = `"\\"${scriptPath.replace(/\\/g, "\\\\")}\\" \\"%V\\""`;

                subCommandList.push({
                    regKey, display, iconPath, command
                });
                subCommands.push(regKey);
            }
        }
        subCommands = `${subCommands.join(";")};`;
        let content = {
            superCommand: {
                project: DIR_PROJECT,
                background: DIR_BACKGROUND,
                icon: path.join(icoDir, TOOL_BOX_ICO).replace(/\\/g, "\\\\")
            },
            subCommands,
            subCommandList
        };
        let addTemp = template.render(fs.readFileSync("./templates/addTemplate.art", "utf-8"), content);
        let deleteTemp = template.render(fs.readFileSync("./templates/deleteTemplate.art", "utf-8"), content);
        let addPath = `${scriptsDir}/toolboxAdd.reg`;
        let deletePath = `${scriptsDir}/toolboxDelete.reg`;
        fs.writeFileSync(addPath,
            addTemp);
        console.log(`${addPath} 生成完成`);
        fs.writeFileSync(deletePath,
            deleteTemp);
        console.log(`${deletePath} 生成完成`);
    } catch (err) {
        console.log(err);
    }

    rd.question("输入任意键结束程序...", () => {
        rd.close();
    });
});
