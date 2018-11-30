import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'
import stylecomponent from '../style/Map/Map.module.css'
import geoJSONmap from '../files/maps/geoJSONmap.json'
import stylemap from '../style/Map/stylemap'
import styletarget from '../style/Targets/styletarget'
import targetsFile from '../files/targets/TARGETS_275.json'

function Map ({ svgMapBounds }) {
  // create map reference
  const mapRef = useRef(null)

  // create map from geoJSON layer
  const maplayer = L.geoJSON(geoJSONmap, {
    style: stylemap
  })

  // create target layer from geoJSON files
  const targetsLayer = L.geoJSON(targetsFile, {
    pointToLayer: function (feature, latlng) {
      return L.circleMarker(latlng, styletarget)
    }
  })

  useEffect(() => {
    mapRef.current = L.map('map', {
      zoomControl: false,
      attributionControl: false,
      zoomSnap: 0.25,
      minZoom: 2.5,
      boxZoom: true
    })

    mapRef.current.addLayer(maplayer)
    mapRef.current.addLayer(targetsLayer)
    mapRef.current.setView([0, 0], 2.5)
    mapRef.current.setMaxBounds(mapRef.current.getBounds())
  }, [])

  return <div id='map' className={stylecomponent.map} style={stylecomponent} />
}

export default Map
