//npm init -y
//npm npm i puppeteer-core
//usage: node script.js /path/to/input.html /path/to/output.pdf


//script.js
const puppeteer = require('puppeteer-core');
const path = require('path');

(async () => {
  
  //let fileinput = __dirname + "/" + process.argv[2];
  let fileinput = "file://" + path.resolve(process.argv[2]);
  console.log(fileinput);
  let fileoutput = process.argv[3];
  //const browser = await puppeteer.launch({args: ['--no-sandbox']});
  const browser = await puppeteer.launch({executablePath: '/usr/bin/google-chrome-stable'});
  const page = await browser.newPage();
  await page.goto(fileinput);
  await page.pdf({path: fileoutput,
    format: "Letter",
    margin: {
      top: '0.8in',
      bottom: '0.8in',
      left: '1.5in',
      right: '1in'
    }, // seems @page in the css overwrite this option
    displayHeaderFooter: true,
    headerTemplate: `<div class='...'></div>`,
    //footerTemplate: `<div style="font-size:10px;text-align:right;width:100%;padding-right:1.0cm;font-family:'Times New Roman', serif;">
    //footerTemplate: `<div style="font-size:10px; font-family:'Times New Roman',serif; width:100%; margin-left:1.4cm; margin-right:1.4cm;margin-bottom:-0.4cm; border: 1px solid #73AD21;">
    //footerTemplate: `<div style="font-size:10px; font-family:'Times New Roman',serif; width:100%; margin-left:1in; margin-right:1in;border: 1px solid #73AD21;">
    // <p style="float:left; text-align:left; width:33.333%;"><span class='date'></span></p>
    /*footerTemplate: `<div style="font-size:10px; font-family:'Times New Roman',serif; width:12.4cm; margin:0 auto;">
    <p style="float:left; text-align:left; width:20%;"></p>
    <p style="float:left; text-align:center; width:60%;">Junli Zhang - Dec 27, 2019</p>
    <p style="float:left; text-align:right; width:20%;"><span class='pageNumber'></span>/<span class='totalPages'></span></p>
    </div>`*/
    //headerTemplate: `<div class='title text center'></div>`,
    // It seems for some reason, footer is scaled by 0.75 of paper size, so if left margin is 1 in, I need to set the footer left margin 0.75 in.
    footerTemplate: `<div style="flex: none;font-size:10px; font-family:'Times New Roman',serif;margin-left:1.125in;border: 1px solid #73AD21;">Left note</div> 
      <div style=" flex: auto;text-align: center;font-size:10px; font-family:'Times New Roman',serif;border: 1px solid #73AD21;">Junli Zhang - Dec 28, 2019</div>
      <div style="flex: none;font-size:10px; font-family:'Times New Roman',serif;margin-right:0.75in;border: 1px solid #73AD21;"><span class='pageNumber'></span>/<span class='totalPages'></span></div>`

  
  });

  await browser.close();
})();
