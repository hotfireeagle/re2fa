import { useEffect } from "react"
import vis from "vis-network/dist/vis-network.esm"
import { options } from "../utils/constant"

const useFA = apiRes => {
  useEffect(() => {
    if (!apiRes) {
      return
    }

    const nodeList = []

    for (let id  = 0; id <= apiRes?.acceptState; id++) {
      const item = { id, label: `s${id}`}
      if (id == apiRes?.startState) {
        item.color = {
          background: "red",
          border: "red",
        }
        item.font = {
          color: "#fff",
        }
      }
      if (id == apiRes?.acceptState) {
        item.color = {
          background: "blue",
          border: "blue",
        }
        item.font = {
          color: "#fff",
        }
      }
      nodeList.push(item)
    }

    // create an array with nodes
    const nodes = new vis.DataSet(nodeList)

    // create an array with edges
    const edges = new vis.DataSet(apiRes.edges)

    // create a network
    const container = document.getElementById("fa")
    const data = { nodes: nodes, edges: edges }

    new vis.Network(container, data, options)
  }, [apiRes])
}

export default useFA