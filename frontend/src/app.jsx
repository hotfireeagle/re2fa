import vis from "vis-network/dist/vis-network.esm"
import { useEffect, useState } from "react"
import request from "./utils/request"

export default function App() {
  const [regexp, setRegexp] = useState("")
  const [apiRes, setApiRes] = useState({})

  useEffect(() => {
    if (!apiRes?.startState) {
      return
    }

    const startNodeId = apiRes.startState.id
    const endNodeId = apiRes.endState.id

    const nodeList = []
    for (let i = startNodeId; i <= endNodeId; i++) {
      const item = { id: i, label: `s${i}`}
      nodeList.push(item)
    }

    // create an array with nodes
    const nodes = new vis.DataSet(nodeList)

    // create an array with edges
    const edges = new vis.DataSet(dfsGenerateEdges(apiRes.startState))

    // create a network
    var container = document.getElementById("fa");
    var data = {
      nodes: nodes,
      edges: edges
    };
    var options = {};
    var network = new vis.Network(container, data, options);
  }, [apiRes])

  const dfsGenerateEdges = startStateObj => {
    const edgeList = []

    const dfs = obj => {
      if (!obj) {
        return
      }
      const { id, transitions } = obj
      const inputSymbols = Object.keys(transitions)

      for (const inputSymbol of inputSymbols) {
        const nextNodeList = transitions[inputSymbol]

        for (const nextNode of nextNodeList) {
          const to = nextNode.id
          const label = inputSymbol == "-1" ? "ε" : inputSymbol
          const item = { from: id, to, label }
          edgeList.push(item)
          dfs(nextNode)
        }
      }
    }

    dfs(startStateObj)

    return edgeList
  }

  const regexpChangeInputHandler = event => {
    setRegexp(event.target.value)
  }

  const generateFAHandler = event => {
    event.preventDefault()
    let finalResult = {}
    request("/api/generateFA", { regexp }).then(result => {
      finalResult = result
    }).finally(() => {
      setApiRes(finalResult)
    })
  }

  return (
    <div className="flex flex-col container bg-gray-50 w-screen max-w-none h-screen">
      <form className="w-screen h50 p-8 bg-blue-800">
        <div className="mx-auto">
          <input
            value={regexp}
            onChange={regexpChangeInputHandler}
            className="border border-transparent focus:outline-none focus:ring-2 focus:ring-purple-600 focus:border-transparent"
          />
          <button
            onClick={generateFAHandler}
            className="bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-600 focus:ring-opacity-50 text-white"
          >
            生成FA描述
          </button>
        </div>
      </form>
      <div id="fa" className="flex-grow bg-white">

      </div>
    </div>
    
  )
}