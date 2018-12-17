import React, { useState } from 'react'
import Map from './components/Map'
import style from './style/App/App.module.css'
import DontReadTheComments from './components/subscriptions'

function App () {
  const svgMapBounds = useState([[-90, -180.4], [89.5, 180]])

  return (
    <div className={style.app}>
      <Map svgMapBounds={svgMapBounds} />
      <DontReadTheComments />
    </div>
  )
}

export default App
