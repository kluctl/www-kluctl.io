const { chromium } = require('playwright');
const { spawn } = require('node:child_process');

const width = 1400
const height = 800

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function kubectl(args) {
    return new Promise((resolve, reject) => {
        const p = spawn('kubectl', args, {});

        let stdout = ""
        p.stdout.on('data', (data) => {
            console.log("kubectl stdout: " + data.toString());
            stdout += data.toString()
        });
        p.stderr.on('data', (data) => {
            console.error("kubectl stderr: " + data.toString());
        });

        p.on("exit", (code) => {
            if (code) {
                reject("process exited with code " + code)
            } else {
                resolve(stdout)
            }
        })
    })
}

/** @ignore */
async function recordWebui(format) {
    const vhs = spawn('vhs', ["../demo.tape"], {
        cwd: "./demo-project"
    });

    const vhsStdout = []
    vhs.stdout.on('data', (data) => {
        console.log("vhs: " + data.toString().trim());
        vhsStdout.push(...data.toString().split("\n"))
    });

    vhs.stderr.on('data', (data) => {
        console.error("vhs stderr: " + data.toString());
    });

    vhs.on('exit', (code) => {
        console.log(`vhs exited with code ${code}`);
    });

    const waitForVhsLine = async (x) => {
        console.log("waiting for vhs line: " + x)
        while (true) {
            const l = vhsStdout.find(l => l.includes(x))
            if (l) {
                console.log("found: " + l)
                break
            }
            await sleep(100)
        }
    }

    const browser = await chromium.launch({
        headless: true,
    });
    const context = await browser.newContext({
        recordVideo: {
            dir: "./recordings",
            size: {
                width: width,
                height: height,
            }
        }
    });

    const page = await context.newPage();
    await page.setViewportSize({width, height})
    page.setDefaultTimeout(60000)

    await page.goto('http://localhost:9090');
    await page.waitForSelector("#kluctl-logo")
    await page.click("#kluctl-logo")
    await page.waitForTimeout(1000);

    const waitForSelector = async (s) => {
        console.log("waiting for " + s)
        await page.waitForSelector(s)
    }

    const waitForFirstCommandResult = async (type, targetName) => {
        const idSuffix = type + "-" + targetName
        await waitForSelector("#firstCommandResult-" + idSuffix)
    }

    const waitForAdditionalCommandResults = async (type, targetName) => {
        const idSuffix = type + "-" + targetName
        await waitForSelector("#additionalCommandResult-" + idSuffix)
    }

    const showFirstCommandResult = async (type, targetName) => {
        const idSuffix = type + "-" + targetName

        // open the command result
        await page.click("#firstCommandResult-" + idSuffix)
        await page.waitForTimeout(2000)

        // click on the changes tab
        await waitForSelector("[id*='T-Changes']")
        await page.click("[id*='T-Changes']")
        await page.waitForTimeout(3000)

        // close the card
        await waitForSelector("#close-card-button")
        await page.click("#close-card-button")
        await page.waitForTimeout(2000)
    }

    // cli

    await waitForFirstCommandResult("cli", "demo-test")
    await showFirstCommandResult("cli", "demo-test")
    await waitForAdditionalCommandResults("cli", "demo-test")
    await showFirstCommandResult("cli", "demo-test")

    // gitops

    await waitForFirstCommandResult("kd", "demo-test")
    await sleep(2000)
    await showFirstCommandResult("kd", "demo-test")

    await waitForVhsLine("git push origin demo-branch")
    await sleep(3000)
    await kubectl(["-n", "kluctl-system", "annotate", "kluctldeployments", "demo", `kluctl.io/request-reconcile='${new Date().toString()}'`])

    await waitForAdditionalCommandResults("kd", "demo-test")
    await showFirstCommandResult("kd", "demo-test")

    await sleep(5000)

    //await recorder.stop();
    await browser.close();
}

recordWebui('./demo-webui.mp4').then(() => {
    console.log('completed');
});
