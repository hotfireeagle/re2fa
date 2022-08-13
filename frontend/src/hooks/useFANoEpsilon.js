import { useEffect } from "react"
import vis from "vis-network/dist/vis-network.esm"
import { options } from "../utils/constant"

const useFANoEpsilon = apiRes => {
  useEffect(() => {
    if (!apiRes) {
      return
    }

    const nodeList = []
    const startState = apiRes?.startState
    const acceptStates = apiRes?.acceptStates

    for (let id of apiRes?.nodes) {
      const item = { id, label: `s${id}`}
      if (startState == id) {
        item.color = {
          background: "red",
          border: "red",
        }
        item.font = {
          color: "#fff",
        }
      }
      if (acceptStates.includes(id)) {
        item.color = {
          background: "blue",
          border: "blue",
        }
        item.font = {
          color: "#fff",
        }
      }
      if (acceptStates.includes(id) && startState == id) {
        item.color = {
          background: "green",
          border: "green",
        }
        item.font = {
          color: "#fff",
        }
      }
      if (id === apiRes?.deadState) {
        item.color = {
          background: "black",
          border: "black",
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
    const newEdges = apiRes.edges.filter(obj => {
      return obj.from != apiRes.deadState
    })
    const edges = new vis.DataSet(newEdges)

    // create a network
    const container = document.getElementById("fa")
    const data = { nodes: nodes, edges: edges }

    new vis.Network(container, data, options)
  }, [apiRes])
}

export default useFANoEpsilon