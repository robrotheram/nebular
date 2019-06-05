import React from 'react';
import {api} from "../utils/api";

export default class Header extends React.Component {
    constructor(...args) {
        super(...args);
        this.state = {user:{username:""}};
    }

    user(){
        api.user().then( data =>{
            this.setState({user:data})
        })
    }
    componentDidMount() {
        this.user();
      } 
    
          
    render() {
      return (
            <header>
                <nav className="navbar navbar-expand-lg navbar-dark bg-primary fixed-top ">
                    <div className="d-flex w-50 order-0">
                        <a className="navbar-brand mr-1" href="/">Nebular</a>
                    </div>
                    <div className="justify-content-center order-2" id="collapsingNavbar">
                    <div className="d-flex justify-content-center h-100">
                    <div className="searchbar">
                    <input className="search_input" type="text" name="search" placeholder="Search..." value={this.props.search} onChange={this.props.onSearch} />
                    <a href="/" className="search_icon"><i className="fa fa-search"></i></a>
                    </div>
                </div>
                    </div>
                    <div className="mt-1 w-50 text-right order-1 order-md-last text-white">{this.state.user.username}</div>
                </nav>
        </header>
      );
    }
  }


