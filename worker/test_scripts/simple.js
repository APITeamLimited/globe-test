import http from 'k6/http';
import ***REMOVED*** sleep ***REMOVED*** from 'k6';

export function contacts() ***REMOVED***
  //const res = http.get('https://test.k6.io/contacts.php', ***REMOVED***
  ////  tags: ***REMOVED*** my_custom_tag: 'contacts' ***REMOVED***,
  ////***REMOVED***);
  //console.log('contacts');
***REMOVED***

export function news() ***REMOVED***
  const startTime = new Date();
  const res = http.get('https://test.k6.io/news.php', ***REMOVED*** tags: ***REMOVED*** my_custom_tag: 'news' ***REMOVED*** ***REMOVED***);
  const endTime = new Date();
  
  console.log(***REMOVED***
    "Yeet": "yeet",
    res,
    startTime,
    endTime,
  ***REMOVED***)

  sleep(1);
***REMOVED***