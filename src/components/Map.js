import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'

// Map outline in form of geojson points
import geoJSONmap from '../files/maps/geoJSONmap.json'

// Style
import stylecomponent from '../style/Map/Map.module.css'
import stylemap from '../style/Map/stylemap'
import styletarget from '../style/Targets/styletarget'
import stylesatellite from '../style/Satellites/stylesatellites'

// Apollo
import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'

// Queries
import gql from 'graphql-tag'
import SatelliteQuery from './../queries/satelliteQuery'
import TargetsQuery from './../queries/targetsQuery'

function Map () {
  const client = new ApolloClient({
    uri: 'http://localhost:3000/graphql',
    link: new HttpLink(),
    cache: new InMemoryCache({
      addTypename: false
    })
  })

  // create map reference
  const mapRef = useRef(null)

  // create map from geoJSON layer
  const maplayer = L.geoJSON(geoJSONmap, {
    style: stylemap
  })

  useEffect(() => {
    client.query({
      query: TargetsQuery
    })
      .then(result => {
        const targetsLayer = L.geoJSON(result.data.targetFeatureCollection, {
          pointToLayer: function (feature, latlng) {
            return L.circleMarker(latlng, styletarget)
          }
        })

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
      })
      .catch(error => console.error(error))
  }, [])

  useEffect(() => {
    const satelliteQuery = client.watchQuery({
      query: SatelliteQuery,
      pollInterval: 1000
    })

    const satelliteSVGstring = "<svg id='satellite' viewBox='0 0 131.3 113.5' className='footer__form__button__svg'><path className='footer__button__svg__path' d='M58.5,76.3l-10.3-8.2l25.2-32l10.3,8.2L58.5,76.3z M99,65.7l-17.5,3.4L82,87.5l31.3,25l17.1-21.8L99,65.7z,M32.3,47.8L32.3,47.8l17.5-3.4L49.4,26L18.1,1L1,22.8L32.3,47.8z M81.5,69.1l-10.7-8.3 M49.7,44.2l9.9,7.7' stroke='#ffffff' strokeWidth='4' /></svg>"

    const satelliteIconUri = encodeURI('data:image/svg+xml,' + satelliteSVGstring).replace('#', '%23')

    const satelliteHere = L.icon({
      iconUrl: satelliteIconUri,
      iconSize: [38, 95]
    })

    const plotSats = (satellites) => {
      satellites.map(satellite => {
        console.log(satellite)
        return L.marker([satellite.latitude, satellite.longitude], { icon: satelliteHere }).addTo(mapRef.current).bindPopup(`${satellite.id} longitude: ${satellite.longitude}`)
      })
    }

    satelliteQuery.subscribe({
      next: (result) => plotSats(result.data.satellites)
    })
  }, [])

  return (
    <div id='map' className={stylecomponent.map} style={stylecomponent} />
  )
}

export default Map
