import http from 'k6/http';
import { sleep } from 'k6';

export function contacts() {
  const res = http.get('https://test.k6.io/contacts.php', {
    tags: { my_custom_tag: 'contacts' },
  });
  console.log('contacts');
  sleep(1);
}

export function news() {
  const res = http.get('https://test.k6.io/news.php', { tags: { my_custom_tag: 'news' } });
  sleep(1);
}