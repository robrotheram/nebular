import React from 'react';
import logo from './logo.svg';
import './App.scss';
import { Dropdown, Modal, Button, Form, Col, OverlayTrigger, Tooltip, Card, ListGroup, Row, Tab, Tabs, Table } from 'react-bootstrap';


import {roles, user} from './data'
const Markdown = require('react-markdown')

class CustomToggle extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e) {
    e.preventDefault();

    this.props.onClick(e);
  }

  render() {
    return (
      <a href="" onClick={this.handleClick} className="dropdownToggle">
        {this.props.children}
      </a>
    );
  }
}


class AddModal extends React.Component {
  render() {
    return (
      <Modal
        {...this.props}
        dialogClassName="modal-900w"
        aria-labelledby="contained-modal-title-vcenter"
        centered
      >
        <Modal.Body>
        <Form>
          <Form.Row>
            <Col>
              <OverlayTrigger key='top' placement='top' overlay={<Tooltip id={`tooltip-top`}>https://github.com</Tooltip>}>
                <Form.Control placeholder="Server" />
              </OverlayTrigger>
            </Col>
            <Col>
              <OverlayTrigger key='top' placement='top' overlay={<Tooltip id={`tooltip-top`}>galaxy</Tooltip>}>
                <Form.Control placeholder="Namespace" />
              </OverlayTrigger>
            </Col>
            <Col>
              <OverlayTrigger key='top' placement='top' overlay={<Tooltip id={`tooltip-top`}>ansible-role</Tooltip>}>
                <Form.Control placeholder="Repository" />
              </OverlayTrigger>
            </Col>
          </Form.Row>
        </Form>
        </Modal.Body>
        <Modal.Footer>
        <Button variant="primary" type="submit" block>Submit</Button>
          <Button onClick={this.props.onHide}>Close</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}


class RepoDetails extends React.Component {
  render() {
    const input = 'How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n How aboutaboutaboutaboutaboutabout some code? \n \n ## How about some code? \n ```js \n var React = \n  \n React.render(\n  <Markdown source="# Your markdown here" />,\n document.getElementById \n);\n```\n some node code'
    const url = this.props.data.Server+'/'+this.props.data.Namespace+'/'+this.props.data.Repo
    return (
      <Modal
        {...this.props}
        dialogClassName="modal-900w"
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
                  <td className="col-2">Repo Url: </td>
                  <td className="col-10"><a href={url}>{url}</a></td>
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
                <h5>Dependancies</h5>
                <ListGroup>
                  <ListGroup.Item>Cras justo odio</ListGroup.Item>
                  <ListGroup.Item>Dapibus ac facilisis in</ListGroup.Item>
                  <ListGroup.Item>Morbi leo risus</ListGroup.Item>
                  <ListGroup.Item>Porta ac consectetur ac</ListGroup.Item>
                  <ListGroup.Item>Vestibulum at eros</ListGroup.Item>
                </ListGroup>
              </Col>
              <Col>
                <h5>Versions</h5>
                <ListGroup>
                  <ListGroup.Item>Cras justo odio</ListGroup.Item>
                  <ListGroup.Item>Dapibus ac facilisis in</ListGroup.Item>
                  <ListGroup.Item>Morbi leo risus</ListGroup.Item>
                  <ListGroup.Item>Porta ac consectetur ac</ListGroup.Item>
                  <ListGroup.Item>Vestibulum at eros</ListGroup.Item>
                </ListGroup>
              </Col>
            </Row>
            </Tab>
            <Tab eventKey="readme" title="ReadMe">
                  <Markdown source={this.props.data.Readme} />
            </Tab>
          </Tabs>
        </Modal.Body>
        <Modal.Footer>
        <div class="tags">
        {this.props.data.Meta.GalaxyInfo.GalaxyTags.map((tag, i) => (
          <a href="#">{tag}</a>
        ))}
    	
    </div>
        </Modal.Footer>
      </Modal>
    );
  }
}

class App extends React.Component {
  constructor(...args) {
    super(...args);

    this.state = { modalShow: false, modalAddShow: false, data:{Meta:{GalaxyInfo:{GalaxyTags:[]}}} };
  }

  render() {
    let modalClose = () => this.setState({ modalShow: false });
    let modalAddClose = () => this.setState({ modalAddShow: false });
    return (
    <div>
      <header>
      <nav class="navbar navbar-expand-lg navbar-dark bg-primary fixed-top ">
        <div class="d-flex w-50 order-0">
            <a class="navbar-brand mr-1" href="#">Nebular</a>
            <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#collapsingNavbar">
                <span class="navbar-toggler-icon"></span>
            </button>
        </div>
        <div class="navbar-collapse collapse justify-content-center order-2" id="collapsingNavbar">
        <div class="d-flex justify-content-center h-100">
        <div class="searchbar">
          <input class="search_input" type="text" name="" placeholder="Search..."/>
          <a href="#" class="search_icon"><i class="fas fa-search"></i></a>
        </div>
      </div>
        </div>
        <div class="mt-1 w-50 text-right order-1 order-md-last">{user.username}</div>
    </nav>
      </header>

    <main role="main" class="container">
      <button type="button" onClick={() => this.setState({ modalAddShow: true })} class="btn btn-primary btn-circle"><i class="fa fa-plus"></i></button>

        <RepoDetails show={this.state.modalShow} onHide={modalClose} data={this.state.data}/>
        <AddModal show={this.state.modalAddShow} onHide={modalAddClose} />
    
  

      <div className="rolelist">
            {roles.map((repo, i) => (
               <div class="card" >
               <div class="card-body">
               <Row>
               <Col>
                <a className="roleName" onClick={() => this.setState({ modalShow: true, data:repo })}>{repo.Namespace}.{repo.Repo}</a>
               </Col>
               <div className="forceRight">
               <Dropdown>
                <Dropdown.Toggle as={CustomToggle} id="dropdown-custom-components">
                <i class="fa fa-ellipsis-v"></i>
                </Dropdown.Toggle>
                    <Dropdown.Menu>
                    <Dropdown.Item href="#/action-1">Delete Role</Dropdown.Item>
                    <Dropdown.Item href="#/action-2">Refresh Role</Dropdown.Item>
                  </Dropdown.Menu>
                </Dropdown> 
               </div>
               </Row>
               </div>
               </div>
              ))}
     </div>

    </main>

    <footer class="footer">
      <div class="container">
        <span class="text-muted">Nebular crated by @robrobotheram</span>
      </div>
    </footer>
  </div>

);
}
}


export default App;
