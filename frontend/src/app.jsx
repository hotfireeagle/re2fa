import { useState } from "react"
import request from "./utils/request"
import useFANoEpsilon from "./hooks/useFANoEpsilon"
import "./index.css"


export default function App() {
  const [regexp, setRegexp] = useState("")
  const [str, setStr] = useState("")
  const [nfaApiRes, setNfaApiRes] = useState(null)
  const [dfaApiRes, setDfaApiRes] = useState(null)

  const successBgColor = "bg-green-200"
  const failedBgColor = "bg-red-200"

  useFANoEpsilon(nfaApiRes)
  useFANoEpsilon(dfaApiRes)

  const regexpChangeInputHandler = event => {
    setRegexp(event.target.value)
  }

  const strChangeHandler = event => {
    setStr(event.target.value)
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

  const faMatchCallback = (id, result) => {
    const dom = document.getElementById(id)
    if (result) {
      dom.classList.remove(failedBgColor)
      dom.classList.add(successBgColor)
    } else {
      dom.classList.remove(successBgColor)
      dom.classList.add(failedBgColor)
    }
  }

  const matchHandler = event => {
    event.preventDefault()
    const urls = [
      { api: "/api/nfaMatch", id: "nfa", },
      { api: "/api/dfaMatch", id: "dfa", },
    ]
    const fetchs = urls.map(obj => {
      let matchResult = false
      const postData = { regexp, text: str }
      return request(obj.api, postData).then(result => {
        console.log(obj.api, result)
        matchResult = result
      }).finally(() => {
        faMatchCallback(obj.id, matchResult)
      })
    })
    return Promise.all(fetchs)
  }

  const resetHandler = event => {
    event.preventDefault()
    const ids = ["nfa", "dfa"]
    ids.forEach(id => {
      const dom = document.getElementById(id)
      dom.classList.remove(failedBgColor)
      dom.classList.remove(successBgColor)
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
            onClick={event => fetchFA(event)}
            className="bg-blue-600 p-2 px-3 rounded-lg mr-6 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white"
          >
            Generate FA
          </button>
        </div>
      </form>
      <div className="flex-grow bg-white flex flex-row">
        <div id="nfa" className="transition-all c1 c11">
        </div>
        <div id="dfa" className="transition-all c1 c12">
        </div>
      </div>
      <form className="w-screen p-5 bg-blue-600 bg-opacity-60">
        <div className="text-center">
          <button
            onClick={event => resetHandler(event)}
            className="bg-blue-600 p-2 px-3 rounded-lg mr-6 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white w90"
          >
            Reset
          </button>
          <input
            placeholder="Enter string"
            value={str}
            onChange={strChangeHandler}
            className="border w-96 max-w-none p-2 px-3 rounded-lg mr-6 border-transparent focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
          />
          <button
            onClick={event => matchHandler(event)}
            className="bg-blue-600 p-2 px-3 rounded-lg mr-6 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50 text-white w90"
          >
            Test Match
          </button>
        </div>
      </form>
    </div>
  )
}