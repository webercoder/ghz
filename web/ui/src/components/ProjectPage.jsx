import React, { Component } from 'react'
import { Pane } from 'evergreen-ui'
import { Provider, Subscribe } from 'unstated'

import ProjectList from './ProjectList'

import ProjectContainer from '../containers/ProjectContainer'

const projectContainer = new ProjectContainer()

export default class ProjectPage extends Component {
  render () {
    return (
      <Provider inject={[projectContainer]}>
        <Subscribe to={[ProjectContainer]}>
          {(projectStore) => (
            <Pane>
              <ProjectList projectStore={projectStore} />
            </Pane>
          )}
        </Subscribe>
      </Provider >
    )
  }
}
