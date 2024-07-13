#!/usr/bin/env node
const { chromium } = require("playwright");

const getArg = (arg, defaultValue) => {
    const commandIndex = process.argv.findIndex((v) => {
        if (v.startsWith(arg)) {
            return v;
        }
    });

    if (commandIndex === -1) {
        return defaultValue;
    }

    if (process.argv[commandIndex+1] === undefined) {
        return defaultValue;
    }

    return process.argv[commandIndex+1]
}

const formatDate = (date) => {
    const year = date.getFullYear().toString();
    let month = `${date.getMonth() + 1}`;
    let day = `${date.getDate()}`;

    if (month.length < 2) month = `0${month}` 
    if (day.length < 2) day = `0${day}`; 

    return [year, month, day].join('-');
}

const getTimetable = async (page) => {
    const now = new Date();
    const day = getArg("-d", formatDate(now));

    await page.goto(`https://belgium.tomorrowland.com/en/line-up/?page=timetable&day=${day}`);
    const buffer = await page.locator('#planby-wrapper').screenshot(); // { path: `timetable_${day}.png` }
    console.log(buffer.toString('base64'));
}

(async () => {
    const browser = await chromium.launch(); // { headless: false }
    const operation = process.argv[2];
    
    if (!operation) {
        throw Error("Specific operation required. Valid operations: timetable, stages.");
    }
    
    const operations = {
        "timetable": getTimetable
    };
    
    if (!operations[operation]) {
        throw Error("Invalid Operation. Valid operations: timetable, stages.")
    }
    
    let page = await browser.newPage();
    
    await page.setViewportSize({ width: 1920, height: 1080 });
    
    await page.goto('https://belgium.tomorrowland.com/en/line-up/');
    const acceptCookiesBtn = page.locator('#CybotCookiebotDialogBodyLevelButtonLevelOptinAllowAll');
    await acceptCookiesBtn.click({ button: 'left' });
    await page.mouse.move(0, 0);

    await page.locator('#CybotCookiebotDialogBodyLevelButtonLevelOptinAllowAll').waitFor({ state: "detached" });

    await operations[operation](page);
    await browser.close();
})();