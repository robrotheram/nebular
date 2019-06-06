import React from 'react';
import { Modal, ListGroup, Row, Col, Tabs, Tab, Table, Card } from 'react-bootstrap';
import  Markdown  from 'react-markdown'
import {api} from "../utils/api";

export default class RoleDetailModal extends React.Component {

  constructor(...args) {
    super(...args);
    this.state = { useSSHURL:false, gitSSHURL:"git@github.com:" }

  }

  settings(){
    api.settings().then( data =>{
      console.log(data)
        this.setState({useSSHURL:data["user_ssh_url"], gitSSHURL:data["git_SSH_URL"]})
    })
  }
  componentDidMount() {
    this.settings();
  } 

  search = (s) => {
    this.props.search(s)
    this.props.onHide()
  }

    render() {
      let meta = this.props.data.Meta;
      let metaType = this.props.data.MetaType;
      const url = this.props.data.Server+'/'+this.props.data.Namespace+'/'+this.props.data.Repo
      let tags = []
      let dependancies = []
      let platforms = []
      console.log(this.props.data)
      if((this.props.data.Meta.GalaxyInfo.GalaxyTags !== undefined)&&(this.props.data.Meta.GalaxyInfo.GalaxyTags !== null)){
        tags = this.props.data.Meta.GalaxyInfo.GalaxyTags
      }
      if(Array.isArray(meta.Dependencies)){
        dependancies = meta.Dependencies;
      }
      if(Array.isArray(this.props.data.Meta.GalaxyInfo.Platforms)){
        platforms = this.props.data.Meta.GalaxyInfo.Platforms
      }
      let ssh_url = url
      if (this.state.useSSHURL){
        ssh_url = this.state.gitSSHURL+this.props.data.Namespace+'/'+this.props.data.Repo
      }
      return (
        <Modal
          {...this.props}
          dialogClassName="modal-900w modal-height"
          aria-labelledby="contained-modal-title-vcenter"
          centered
        >
         <Modal.Header closeButton style={{paddingBottom:"35px"}}>
            <Modal.Title>{this.props.data.Repo}</Modal.Title>
            <p style={{position:'absolute', top: "50px"}}>{this.props.data.Meta.GalaxyInfo.Description}</p>
          </Modal.Header>
          <Modal.Body>
            <Tabs defaultActiveKey="details" id="uncontrolled-tab-example">
              <Tab eventKey="details" title="Details">
                <br/>
                <Row>
                  <Col>
                    <Card>
                    <Table>
                      <tbody>
                        <tr>
                          <td>Repo Url: </td>
                          <td><a href={url}>{url}</a></td>
                        </tr>
                        <tr>
                          <td>Minimum Ansible Version: </td>
                          <td>{this.props.data.Meta.GalaxyInfo.MinAnsibleVersion}</td>
                        </tr>
                        <tr>
                          <td>Author</td>
                          <td><a href="#/" onClick={() => this.search(this.props.data.Meta.GalaxyInfo.Author)}>{this.props.data.Meta.GalaxyInfo.Author}</a></td>
                        </tr>
                        <tr>
                          <td>Company</td>
                          <td>{this.props.data.Meta.GalaxyInfo.Company}</td>
                        </tr>
                        <tr>
                          <td>License</td>
                          <td>{this.props.data.Meta.GalaxyInfo.License}</td>
                        </tr>
                      </tbody>
                    </Table>
                    </Card>
                  </Col>
                  <Col>
                  <h5>Installing:</h5>
                  To use this role copy the following into your requirements.yml for your playbook
                  <br/><br/>
                  <pre>
                    <code>
{`- name: `+this.props.data.Repo+`
  src: git+`+ssh_url+`
  scm: git
`}
                      </code>
                  </pre>
                  </Col>
                </Row>
              <hr/>
              <Row>
                <Col>
                  <h5>Dependancies:</h5>
                  <ListGroup>         
                  {dependancies.map((dependency, i) => {
                    if (metaType === "COMPLEX"){
                      return (<ListGroup.Item key={i}><a href={dependency.Src}>{dependency.Src}</a></ListGroup.Item>)
                    } else {
                      return (<ListGroup.Item key={i}>{dependency}</ListGroup.Item>)
                    }
                  })}
                  </ListGroup>
                </Col>
                <Col>
                  <h5>Versions:</h5>
                  <ListGroup>
                  {platforms.map((dependency, i) => (
                    <ListGroup.Item key={i}>
                      {dependency.Name}: 
                      <span className="version">
                      {dependency.Versions.map((v, j) => (
                        <i key={j}> {v} </i>
                      ))}
                      </span>
                    </ListGroup.Item>
                  ), this)}
                  </ListGroup>
                </Col>
              </Row>
              </Tab>
              <Tab eventKey="readme" title="ReadMe" className="readme">
                    <Markdown source={this.props.data.Readme} />
              </Tab>
            </Tabs>
          </Modal.Body>
          <Modal.Footer>
          <div className="tags">
            {tags.map((tag, i) => (
              <i key={i} onClick={() => this.search(tag)}>{tag}</i>
            ),this)}
          </div>
          </Modal.Footer>
        </Modal>
      );
    }
  }