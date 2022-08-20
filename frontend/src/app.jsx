import { useState } from "react"
import request from "./utils/request"
import useFANoEpsilon from "./hooks/useFANoEpsilon"
import "./index.css"

// TODO: no more tailwind

export default function App() {
  const [regexp, setRegexp] = useState("")
  const [nfaApiRes, setNfaApiRes] = useState(null)
  const [dfaApiRes, setDfaApiRes] = useState(null)

  useFANoEpsilon(nfaApiRes)
  useFANoEpsilon(dfaApiRes)

  const regexpChangeInputHandler = event => {
    setRegexp(event.target.value)
  }

  const fetchFA = event => {
    event.preventDefault()
    const urls = [
      { api: "/api/generateFA", id: "nfa", cb: setNfaApiRes, },
      { api: "/api/generateDFA", id: "dfa", cb: setDfaApiRes, },
    ]
    const requests = urls.map(obj => {
      let finalResult = null
      return request(obj.api, { regexp }).then(result => {
        finalResult = result
      }).finally(() => {
        obj.cb({ ui: finalResult, id: obj.id })
      })
    })
    return Promise.all(requests)
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
            onClick={event => fetchFA(event)}
            className="bg-blue-600 p-2 px-3 rounded-lg mr-6 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white"
          >
            Generate FA
          </button>
        </div>
      </form>
      <div className="flex-grow bg-white flex flex-row">
        <div id="nfa" className="c1 c11">
        </div>
        <div id="dfa" className="c1 c12">
        </div>
      </div>
    </div>
  )
}