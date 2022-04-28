import http from "k6/http";
import ***REMOVED*** sleep ***REMOVED*** from "k6";

let accessToken = "YOUR_GITHUB_ACCESS_TOKEN";

export default function() ***REMOVED***

  let query = `
    query FindFirstIssue ***REMOVED***
      repository(owner:"grafana", name:"k6") ***REMOVED***
        issues(first:1) ***REMOVED***
          edges ***REMOVED***
            node ***REMOVED***
              id
              number
              title
            ***REMOVED***
          ***REMOVED***
        ***REMOVED***
      ***REMOVED***
    ***REMOVED***`;

  let headers = ***REMOVED***
    'Authorization': `Bearer $***REMOVED***accessToken***REMOVED***`,
    "Content-Type": "application/json"
  ***REMOVED***;

  let res = http.post("https://api.github.com/graphql",
    JSON.stringify(***REMOVED*** query: query ***REMOVED***),
    ***REMOVED***headers: headers***REMOVED***
  );

  if (res.status === 200) ***REMOVED***
    console.log(JSON.stringify(res.body));
    let body = JSON.parse(res.body);
    let issue = body.data.repository.issues.edges[0].node;
    console.log(issue.id, issue.number, issue.title);

    let mutation = `
      mutation AddReactionToIssue ***REMOVED***
        addReaction(input:***REMOVED***subjectId:"$***REMOVED***issue.id***REMOVED***",content:HOORAY***REMOVED***) ***REMOVED***
          reaction ***REMOVED***
            content
          ***REMOVED***
          subject ***REMOVED***
            id
          ***REMOVED***
        ***REMOVED***
    ***REMOVED***`;

    res = http.post("https://api.github.com/graphql",
      JSON.stringify(***REMOVED***query: mutation***REMOVED***),
      ***REMOVED***headers: headers***REMOVED***
    );
  ***REMOVED***
  sleep(0.3);
***REMOVED***
