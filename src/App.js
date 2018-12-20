import React from 'react'
import Map from './components/Map'
import style from './style/App/App.module.css'
import { ApolloClient } from 'apollo-client'
import { ApolloProvider } from 'react-apollo'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'
// import TargetsLayer from './components/TargetsQuery'

function App () {

  const client = new ApolloClient({
    uri: 'http://localhost:3000/graphql',
    link: new HttpLink(),
    cache: new InMemoryCache()
  })

  return (
    <ApolloProvider client={client} >
      <div className={style.app}>
        <Map />
      </div>
    </ApolloProvider>

  )
}

export default App
