import gql from 'graphql-tag'

const TargetsQuery = gql`{
  targetFeatureCollection{
    type,
    features {
      type,
      geometry {
        type,
        coordinates,
      },
      properties {
        shortName,
      },
    }
  }
}`

export default TargetsQuery
