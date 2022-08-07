import vis from "vis-network/dist/vis-network.esm"
import { useEffect, useState } from "react"
import request from "./utils/request"

export default function App() {
  const [regexp, setRegexp] = useState("")
  const [apiRes, setApiRes] = useState({})

  useEffect(() => {
    if (!apiRes?.nodes) {
      return
    }

    const nodeList = []
    for (const id of apiRes.nodes) {
      const item = { id, label: `s${id}`}
      nodeList.push(item)
    }

    // create an array with nodes
    const nodes = new vis.DataSet(nodeList)

    // create an array with edges
    const edges = new vis.DataSet(apiRes.edges)

    // create a network
    const container = document.getElementById("fa")
    const data = { nodes: nodes, edges: edges }

    const options = {
      edges: {
        arrows: {
          to: {
            enabled: true,
            scaleFactor: 1,
            type: "arrow"
          },
          from: {
            enabled: false,
            scaleFactor: 1,
            type: "arrow"
          }
        },
        arrowStrikethrough: true,
        chosen: true,
        // transition line color
        color: {
          color: "#848484",
          highlight: "#848484",
          hover: "#848484",
          inherit: "from",
          opacity: 1.0
        },
        dashes: false,
        hoverWidth: 1.5,
        labelHighlightBold: true,
        physics: true,
      },
      layout: {
        randomSeed: 0.1,
      }
    }
    new vis.Network(container, data, options)
  }, [apiRes])

  const regexpChangeInputHandler = event => {
    setRegexp(event.target.value)
  }

  const generateFAHandler = event => {
    event.preventDefault()
    let finalResult = {}
    return request("/api/generateFA", { regexp }).then(result => {
      finalResult = result
    }).finally(() => {
      setApiRes(finalResult)
    })
  }

  return (
    <div className="flex flex-col container bg-gray-50 w-screen max-w-none h-screen">
      <form className="w-screen p-5 bg-blue-600 bg-opacity-60">
        <div className="text-center">
          <input
            placeholder="Enter regexp"
            value={regexp}
            onChange={regexpChangeInputHandler}
            className="border w-96 max-w-none p-2 px-3 rounded-lg mr-6 border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
          />
          <button
            onClick={generateFAHandler}
            className="bg-blue-600 p-2 px-3 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white"
          >
            Generate FA
          </button>
        </div>
      </form>
      <div id="fa" className="flex-grow bg-white" />
    </div>
  )
}