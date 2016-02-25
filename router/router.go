package router

import (
	"github.com/prime/handler"
	"gopkg.in/macaron.v1"
)

func SetRouters(m *macaron.Macaron) {
	m.Group("/repos/:owner/:repo", func() {
		//Git Data
/*		m.Group("/git", func() {
			//blobs
		    m.Get("/blobs/:sha", handler.GetBlobsHandler)  
		    m.Post("/blobs", handler.CreateBlobsHandler)
			//commits		    
		    m.Get("/commits/:sha", handler.GetCommitsHandler)
		    m.Post("/commits", handler.CreateCommitsHandler)
			//Reference		    
		    m.Get("/refs/:ref", handler.GetSingleRefHandler)
		    m.Get("/refs", handler.GetAllRefsHandler)
		    m.Post("/refs", handler.CreateSingleRefHandler)
		    m.Patch("/refs/:ref", handler.UpdateSingleRefHandler)
			//Tag		    
		    m.Get("/tags/:sha", handler.GetTagHandler)
		    m.Post("/tags", handler.CreateTagHandler)	
			//Tree		    
		    m.Get("/trees/:sha", handler.GetTreeHandler)
		    m.Get("/trees/:sha?recursive=1", handler.GetTreeRecursivelyHandler)
		    m.Post("/trees", handler.CreateTreeHandler)		    
		})
*/
        //Pull Request
		m.Get("/pulls", handler.ListPullRequestHandler)  
		m.Get("/pulls/:number", handler.GetSinglePullRequestHandler)  
		m.Post("/pulls", handler.CreateNewPullHandler)
		m.Patch("/pulls/:number", handler.UpdateSinglePullHandler)  
		m.Get("/pulls/:number/commits", handler.ListCommitsOnPullHandler)  
		m.Get("/pulls/:number/files", handler.ListFilesOnPullHandler)  
		m.Get("/pulls/:number/merge", handler.IfPullMergedHandler)
		m.Put("/pulls/:number/merge", handler.MergeNewPullHandler)  
	})
    
       //Issues
    
}
