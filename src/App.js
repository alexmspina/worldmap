import React, { useState } from 'react'
import Map from './components/Map'
import style from './style/App/App.module.css'

function App () {
  const [markerPosition, setMarkerPosition] = useState({
    lat: -60,
    lng: -170
  })
  const { lat, lng } = markerPosition

  function moveMarker () {
    setMarkerPosition({
      lat: lat + 0.0001,
      lng: lng + 0.0001
    })
  }

  return (
    <div className={style.app}>
      <Map markerPosition={markerPosition} />
    </div>
  )
}

export default App
