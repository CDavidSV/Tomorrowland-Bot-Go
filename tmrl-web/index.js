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

(async () => {
    const now = new Date();
    const day = getArg("-d", formatDate(now));

    const browser = await chromium.launch();

    let page = await browser.newPage();
    await page.setViewportSize({ width: 2560, height: 1440 });
    await page.goto(`https://belgium.tomorrowland.com/en/line-up/?page=timetable&day=${day}`);

    await page.locator('#CybotCookiebotDialogBodyLevelButtonLevelOptinAllowAll').click({ button: 'left' });
    await page.mouse.move(0, 0);

    setTimeout(async () => {
        const buffer = await page.locator('#planby-wrapper').screenshot(); // { path: `timetable_${day}.png` }
        console.log(buffer.toString('base64'));

        await browser.close();
    }, 300);
})();