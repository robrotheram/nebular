import React from 'react';
import {api} from "./utils/api";
import './App.css';

import Header from './components/header'
import RoleDetailModal from './components/roleDetail'
import AddRoleModal from './components/addRole'

import { Dropdown, Col, Row} from 'react-bootstrap';




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

class Nebular extends React.Component {
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

  
  componentDidMount() {
    this.getRoles();
    
  } 

  handleSearch = (event) => {
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
    return (
    <div>
      <Header search={this.state.search} onSearch={this.handleSearch}/>
      <RoleDetailModal show={this.state.modalShow} onHide={modalClose} data={this.state.data}/>
      <AddRoleModal show={this.state.modalAddShow} onHide={modalAddClose} refresh={this.getRoles}/>
      <button type="button" onClick={() => this.setState({ modalAddShow: true })} className="btn btn-primary btn-circle"><i className="fa fa-plus"></i></button>
      <main role="main" className="container">
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


export default Nebular;
