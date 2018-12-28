import gql from 'graphql-tag'

const CatseyesQuery = gql`{
  catseyeFeatureCollection{
    type,
    features {
      type,
      geometry {
        type,
        coordinates,
      },
      properties {
        subregion,
      },
    }
  }
}`

export default CatseyesQuery
