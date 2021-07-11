package blockQuery

import (
	"encoding/json"
	"frontier/api"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

//Controller model and its parameters
type Controller struct {
	log *zap.SugaredLogger
	mgr *Manager
}

// NewBlockQueryController Initializing the required parameters to return the controller instance
func NewBlockQueryController(logger *zap.SugaredLogger, manager *Manager) *Controller{
  return &Controller{log: logger,mgr: manager}
}

// GetPaths Initializing the BaseURL for the block_query query
func (controller *Controller) GetPaths() []api.HTTPRoute {
	baseURL := "/block_query"
	routes := []api.HTTPRoute{
		{Method: "GET", Path: baseURL},
	}
	return routes
}

/* The GET method handles the incoming HTTP request and returns the appropriate response
   function calls the manager GetRequiredBlocks method and writes the response into output.json
 */
func (controller *Controller) GET(w http.ResponseWriter, r *http.Request) {
	res := controller.mgr.GetRequiredBlocks()
	file, _ := json.MarshalIndent(res, "", " ")
	_ = ioutil.WriteFile("output.json", file, 0644)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	controller.log.Infof("Total No records Fetched %d",len(res.Result))
	//w.Header().Set("Content-Length", strconv.Itoa(len(res.Result)))
	_ = json.NewEncoder(w).Encode(res)
}

func (controller *Controller) LIST(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}
