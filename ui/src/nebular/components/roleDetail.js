import React from 'react';
import { Modal, ListGroup, Row, Col, Tabs, Tab, Table } from 'react-bootstrap';
import  Markdown  from 'react-markdown'

export default class RoleDetailModal extends React.Component {
    render() {
      let meta = this.props.data.Meta;
      let metaType = this.props.data.MetaType;
      const url = this.props.data.Server+'/'+this.props.data.Namespace+'/'+this.props.data.Repo
      let tags = []
      console.log(this.props.data)
      if((this.props.data.Meta.GalaxyInfo.GalaxyTags !== undefined)&&(this.props.data.Meta.GalaxyInfo.GalaxyTags !== null)){
        tags = this.props.data.Meta.GalaxyInfo.GalaxyTags
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
                    <td>{this.props.data.Meta.GalaxyInfo.Author}</td>
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
              <hr/>
              <Row>
                <Col>
                  <h5>Dependancies:</h5>
                  <ListGroup>         
                  {meta.Dependencies.map((dependency, i) => {
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
                  {this.props.data.Meta.GalaxyInfo.Platforms.map((dependency, i) => (
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
              <i key={i}>{tag}</i>
            ))}
          </div>
          </Modal.Footer>
        </Modal>
      );
    }
  }