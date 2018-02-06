package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"log"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/graphql_server"
	"github.com/vulcanize/vulcanizedb/utils"
)

// startGraphqlCmd represents the startGraphql command
var graphqlCmd = &cobra.Command{
	Use:   "graphql",
	Short: "Starts Vulcanize graphql server",
	Long: `Starts vulcanize graphql server
and usage of using your command. For example:

graphql --port 9090 --host localhost
`,
	Run: func(cmd *cobra.Command, args []string) {
		schema := parseSchema()
		serve(schema)
	},
}

func init() {
	var (
		port int
		host string
	)
	rootCmd.AddCommand(graphqlCmd)

	syncCmd.Flags().IntVar(&port, "port", 9090, "graphql: port")
	syncCmd.Flags().StringVar(&host, "host", "localhost", "graphql: host")

}

func parseSchema() *graphql.Schema {

	blockchain := geth.NewBlockchain(ipc)
	repository := utils.LoadPostgres(databaseConfig, blockchain.Node())
	schema := graphql.MustParseSchema(graphql_server.Schema, graphql_server.NewResolver(repository))
	return schema

}

func serve(schema *graphql.Schema) {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))
	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":9090", nil))
}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
