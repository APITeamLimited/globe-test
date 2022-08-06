import http from 'k6/http';
import ***REMOVED*** sleep ***REMOVED*** from 'k6';

export const options = ***REMOVED***
    discardResponseBodies: true,
    scenarios: ***REMOVED***
      contacts: ***REMOVED***
        executor: 'constant-vus',
        //vus: 10000,
        //duration: '60s',
      ***REMOVED***,
    ***REMOVED***,
  ***REMOVED***;

export default function () ***REMOVED***
    const randomPage = Math.floor(Math.random() * 100);

  //http.get(`https://api-staging.inteliscan.app/content/scrapes-analyst/?page=$***REMOVED***randomPage***REMOVED***`, ***REMOVED***
    http.get(`https://api-staging.inteliscan.app/ping/`, ***REMOVED***
    headers: ***REMOVED***
        // Add auth token here
        // /Authorization: 'Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNjU5NzkzMjI0LCJqdGkiOiI5YTJkNTE0OTg3ZDI0ZTYxYmViNWZkYWZlM2IzYjhhYSIsInVzZXJfaWQiOiIyNGJiNmU1NS01ODM5LTQxZWMtOTFiMi1iMDJmOTJkODNhODgiLCJwZXJtaXNzaW9ucyI6WyJHb2QiXSwicm9sZSI6IlN1cGVyIEFkbWluIiwiYXVkIjoiSW50ZWxpc2NhbiIsImlzcyI6IkludGVsaXNjYW4ifQ.vqDqbcut4l7rm2hqBy6PRM63yVj7f7S2whUI7UQQCZOUXR8AxiNVlkFQ1fmhRn32MIa-HpApOKCFY4aaSXdDjeS43QeadKMeNbJlvq5PATv_629i-_KUeIsziyGU959R7nMkS9trosW-uyLmiDa5HvkHbGmS-kiby7RT9TRbq5Xsias9suGoaIxgpta3k1s9rVunYwsLcu2MPrXqVfgGO25TKyGqzd9HGhAlD13DhTWiLkmhXtGq91POKvi06vwzJF6EgRqfF4Kv_HE2AGF0gK9fs6QXu1OKuieU2KV7EGenyako7XmEZgbuI498tymZfXxlzfpW_3L5X1Z2ly6OqyMnTrvNXFC2hx7LgpFGu6Yy56_gR5j0Ky_t4p4e26bmrUM3Pqb_Lb1BNZs2RWjDs6rpuZ_NYN9hJurrJG5v_JeV9Bj2c83orzt3mCky5AmlEENFrt0ONQ-t_o0Md7_KXXPJG0QaQggcTHsdjote6uq-K2d7pXtuRWMaPMbzPAM74GmLWMP7ACvbShuw4iLpfUF6JJ4JjykKWE5l5pA7FHMVrZAfvoPGSwANyUmWHFutBbYRMvmB0c5jsh3WD2b8ulk7Lh4YfpBl0hkj5ee2DR_R_dPXqXrbOA3Km4lYzibN3pCOxhC0ro5KTyMmMs1FREXstULYiPQ9pxfNglvkxk4',
    ***REMOVED***,
  ***REMOVED***)
    sleep(10)
***REMOVED***
