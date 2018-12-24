import React from 'react'
import WorldMap from './components/Map'
import style from './style/App/App.module.css'

function App () {
  return (
    <div className={style.app}>
      <WorldMap id='map' />
    </div>

  )
}

export default App
