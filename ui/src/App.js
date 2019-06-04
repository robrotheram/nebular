import React from 'react';
import './App.css';
import { Dropdown, Modal, Button, Form, Col, OverlayTrigger, Tooltip, ListGroup, Row, Tab, Tabs, Table } from 'react-bootstrap';

import {api} from "./api";

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
      <a href="/" onClick={this.handleClick} className="dropdownToggle">
        {this.props.children}
      </a>
    );
  }
}


class AddModal extends React.Component {

  constructor(...args) {
    super(...args);
    this.default = { Server:"https://github.com", Namespace:"", Repository:"test" }
    this.state = this.default;
    this.handleChange = this.handleChange.bind(this);
    this.close = this.close.bind(this);
    this.addModel = this.addModel.bind(this);
  }

  handleChange(event) {
    this.setState({ [event.target.name]: event.target.value });
  }

  close(){
    this.setState(this.default)
    this.props.onHide()
  }

  addModel(){
    let _this = this;
    api.create({
      "Server": this.state.Server,
      "Namespace": this.state.Namespace,
      "Repo": this.state.Repository
    }).then(data =>{ 
      console.log(data)
      _this.props.refresh()
      _this.props.onHide()

    })
  }

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
                <Form.Control placeholder="https://github.com" name="Server" value={this.state.Server} onChange={this.handleChange} />
              </OverlayTrigger>
            </Col>
            <Col>
              <OverlayTrigger key='top' placement='top' overlay={<Tooltip id={`tooltip-top`}>galaxy</Tooltip>}>
                <Form.Control placeholder="Namespace" name="Namespace" value={this.state.Namespace} onChange={this.handleChange} />
              </OverlayTrigger>
            </Col>
            <Col>
              <OverlayTrigger key='top' placement='top' overlay={<Tooltip id={`tooltip-top`}>ansible-role</Tooltip>}>
                <Form.Control placeholder="Repository" name="Repository" value={this.state.Repository} onChange={this.handleChange} />
              </OverlayTrigger>
            </Col>
          </Form.Row>
        </Form>
        </Modal.Body>
        <Modal.Footer>
        <Button variant="primary" type="submit" onClick={this.addModel} block>Submit</Button>
          <Button onClick={this.close}>Close</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}


class RepoDetails extends React.Component {
  render() {
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
                <h5>Dependancies:</h5>
                <ListGroup>
                {this.props.data.Meta.Dependencies.map((dependency, i) => (
                  <ListGroup.Item key={i}>{dependency}</ListGroup.Item>
                ), this)}
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

class App extends React.Component {
  constructor(...args) {
    super(...args);

    this.state = { 
      roles: [], 
      search:"", 
      user:{username:""},
      modalShow: false, 
      modalAddShow: false, 
      data:{
        Meta:{
          Dependencies: [],
          GalaxyInfo:{
            Platforms: [],
            GalaxyTags:[]
          }
        }
      }
    };
    this.handleChange = this.handleChange.bind(this);
    this.getRoles = this.getRoles.bind(this);
  }

  getRoles(){
    api.getAll().then( data =>{
      this.setState({roles:data})
    }
    )
  }

  search(){
    api.search(this.state.search).then( data =>{
      this.setState({roles:data})
    })
  }

  user(){
    api.user().then( data =>{
      this.setState({user:data})
    })
  }
  componentDidMount() {
    this.getRoles();
    this.user();
  } 

  handleChange(event) {
    if(event === undefined){return}
    this.setState({search: event.target.value});
    console.log(event.target.value)
    if(event.target.value===""){
      api.getAll().then( data =>{
        this.setState({roles:data})
      })
    }else{
      api.search(event.target.value).then( data =>{
        this.setState({roles:data})
      })
    }
    
  }

  delete(id){
    let _this = this;
    console.log("deleting", id)
    api.delete(id).then(data =>{
      console.log(data)
      _this.getRoles()
    }) 
  }

  update(id){
    let _this = this;
    console.log("updating", id)
    api.update(id).then(data =>{
      console.log(data)
      _this.getRoles()
    }) 
  }

  render() {
    let modalClose = () => this.setState({ modalShow: false });
    let modalAddClose = () => this.setState({ modalAddShow: false });
    let user = this.state.user
    return (
    <div>
      <header>
      <nav className="navbar navbar-expand-lg navbar-dark bg-primary fixed-top ">
        <div className="d-flex w-50 order-0">
            <a className="navbar-brand mr-1" href="/">Nebular</a>
        </div>
        <div className="justify-content-center order-2" id="collapsingNavbar">
        <div className="d-flex justify-content-center h-100">
        <div className="searchbar">
          <input className="search_input" type="text" name="search" placeholder="Search..." value={this.state.search} onChange={this.handleChange} />
          <a href="/" className="search_icon"><i className="fa fa-search"></i></a>
        </div>
      </div>
        </div>
        <div className="mt-1 w-50 text-right order-1 order-md-last text-white">{user.username}</div>
    </nav>
      </header>

    <main role="main" className="container">
      <button type="button" onClick={() => this.setState({ modalAddShow: true })} className="btn btn-primary btn-circle"><i className="fa fa-plus"></i></button>

        <RepoDetails show={this.state.modalShow} onHide={modalClose} data={this.state.data}/>
        <AddModal show={this.state.modalAddShow} onHide={modalAddClose} refresh={this.getRoles}/>
    
  

      <div className="rolelist">
            {this.state.roles.map((repo, i) => (
               <div className="card" key={i} >
               <div className="card-body">
               <Row>
               <Col>
                <a className="roleName" onClick={() => this.setState({ modalShow: true, data:repo })}>
                {repo.Namespace}.{repo.Repo} | {repo.Meta.GalaxyInfo.Description}
                </a>
               </Col>
               <div className="forceRight">
               <Dropdown>
                <Dropdown.Toggle as={CustomToggle} id="dropdown-custom-components">
                <i className="fa fa-ellipsis-v"></i>
                </Dropdown.Toggle>
                    <Dropdown.Menu>
                    <Dropdown.Item onClick={() => this.delete(repo.ID)}>Delete Role</Dropdown.Item>
                    <Dropdown.Item onClick={() => this.update(repo.ID)}>Refresh Role</Dropdown.Item>
                  </Dropdown.Menu>
                </Dropdown> 
               </div>
               </Row>
               </div>
               </div>
              ), this)}
     </div>

    </main>

    <footer className="footer">
      <div className="container">
        <span className="text-muted">Nebular crated by @robrobotheram</span>
      </div>
    </footer>
  </div>

);
}
}


export default App;
