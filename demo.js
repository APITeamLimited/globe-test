import http from 'k6/http';
import ***REMOVED*** sleep ***REMOVED*** from 'k6';

export function contacts() ***REMOVED***
  const res = http.get('https://test.k6.io/contacts.php', ***REMOVED***
    tags: ***REMOVED*** my_custom_tag: 'contacts' ***REMOVED***,
  ***REMOVED***);
  console.log('contacts');
  sleep(1);
***REMOVED***

export function news() ***REMOVED***
  const res = http.get('https://test.k6.io/news.php', ***REMOVED*** tags: ***REMOVED*** my_custom_tag: 'news' ***REMOVED*** ***REMOVED***);
  sleep(1);
***REMOVED***