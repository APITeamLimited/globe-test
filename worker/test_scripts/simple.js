import http from 'k6/http';
import { sleep } from 'k6';

export function contacts() {
  //const res = http.get('https://test.k6.io/contacts.php', {
  ////  tags: { my_custom_tag: 'contacts' },
  ////});
  //console.log('contacts');
}

export function news() {
  const startTime = new Date();
  const res = http.get('https://test.k6.io/news.php', { tags: { my_custom_tag: 'news' } });
  const endTime = new Date();
  
  console.log({
    "Yeet": "yeet",
    res,
    startTime,
    endTime,
  })

  sleep(1);
}