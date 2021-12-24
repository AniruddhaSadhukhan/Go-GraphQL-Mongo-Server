package gqlhandler

import (
	"go-graphql-mongo-server/logger"
	"net/http"
)

func GraphiqlHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(graphiQlPage); err != nil {
		logger.Log.Error("Error %v with request %v", err, r)
	}
}

var graphiQlPage = []byte(`
<!DOCTYPE html>
<html>
  <head>
    <title>GraphiQL</title>
    <link href="https://unpkg.com/graphiql/graphiql.min.css" rel="stylesheet" />
  </head>
  <body style="margin: 0;">
    <div id="graphiql" style="height: 100vh;"></div>

    <script
      crossorigin
      src="https://unpkg.com/react/umd/react.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/react-dom/umd/react-dom.production.min.js"
    ></script>
    <script
      crossorigin
      src="https://unpkg.com/graphiql/graphiql.min.js"
    ></script>

    <script>
      const fetcher = GraphiQL.createFetcher({ url: '/api/v1/graphql' });

      ReactDOM.render(
        React.createElement(GraphiQL, { 
			fetcher: fetcher,
			docExplorerOpen: true,
			shouldPersistHeaders: true
		 }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>

`)
