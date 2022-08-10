import launcher from 'k6/x/browser';

export default function () ***REMOVED***
  const browser = launcher.launch('chromium', ***REMOVED*** headless: false ***REMOVED***);
  const context = browser.newContext();
  const page = context.newPage();
  page.goto('http://whatsmyuseragent.org/');
  page.screenshot(***REMOVED*** path: `example-chromium.png` ***REMOVED***);
  console.log("Done");
  page.close();
  browser.close();
***REMOVED***