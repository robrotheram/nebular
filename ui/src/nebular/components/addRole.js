import React from 'react';
import {api} from "../utils/api";
import { Modal, Button, Form, Col, OverlayTrigger, Tooltip } from 'react-bootstrap';

export default class AddRoleModal extends React.Component {

    constructor(...args) {
      super(...args);
      this.default = { Server:"https://github.com", Namespace:"", Repository:"" }
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
    settings(){
      api.settings().then( data =>{
        console.log(data)
          this.setState({Server:data["git_server"], Namespace:data["git_namespace"]})
      })
    }
    componentDidMount() {
      this.settings();
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