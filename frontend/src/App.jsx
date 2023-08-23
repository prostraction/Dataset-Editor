import * as React from 'react';
import './App.css';
import { SetDirectoryDialog } from '../wailsjs/go/main/App';
import { StartMergeProcess } from '../wailsjs/go/main/App';
import { StartCropProcess } from '../wailsjs/go/main/App';

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isMergeEnabled: false,
      dir1: 'Files directory', dir1Selected: false,
      dir2: 'Second files directory', dir2Selected: false,
      dir3: 'Result Directory', dir3Selected: false,
    };

    this.SetDirectory1 = this.SetDirectory1.bind(this);
    this.SetDirectory2 = this.SetDirectory2.bind(this);
    this.SetDirectory3 = this.SetDirectory3.bind(this);
    this.ChangeMode = this.ChangeMode.bind(this);
    this.Progress = this.Progress.bind(this);
  }

  SetDirectory1 = () => {
    SetDirectoryDialog().then((result) => {
      if (result !== "Dialog cancelled") {
        this.setState({ dir1: result, dir1Selected: true});
        document.body.querySelector('#selectdir1').style.color = '#FFFFFF'
      }
    });
  };
  SetDirectory2 = () => {
    SetDirectoryDialog().then((result) => {
      if (result !== "Dialog cancelled") {
        this.setState({ dir2: result, dir2Selected: true});
        document.body.querySelector('#selectdir2').style.color = '#FFFFFF'
      }
    });
  };
  SetDirectory3 = () => {
    SetDirectoryDialog().then((result) => {
      if (result !== "Dialog cancelled") {
        this.setState({ dir3: result, dir3Selected: result !== true});
        document.body.querySelector('#selectdir3').style.color = '#FFFFFF'
      }
    });
  };

  ChangeMode = () => {
    let flag = this.state.isMergeEnabled;
    let dirText = this.state.dir1;
    if (!this.state.dir1Selected) {
      dirText = flag ? 'Files directory' : 'First files directory';
    }
    document.body.querySelector('#selectdir1').style.color = '#FFFFFF'
    document.body.querySelector('#selectdir2').style.color = '#FFFFFF'
    document.body.querySelector('#selectdir3').style.color = '#FFFFFF'
    this.setState({ isMergeEnabled: !flag, dir1: dirText });
  };

  Progress = () => {
    let ready = true;
    if (!this.state.dir1Selected) {document.body.querySelector('#selectdir1').style.color = '#D70040'; ready = false;}
    if (!this.state.dir2Selected && this.state.isMergeEnabled) {document.body.querySelector('#selectdir2').style.color = '#D70040'; ready = false;}
    if (!this.state.dir3Selected) {document.body.querySelector('#selectdir3').style.color = '#D70040'; ready = false;}
    
    let x = document.getElementById('xInput').value;
    let y = document.getElementById('yInput').value;
    x = parseInt(x); 
    y = parseInt(y);

    let queryX = document.body.querySelector('#xInput'); 
    if (isNaN(x)) {
      queryX.style.color = '#D70040'; ready = false;
    } else {
      queryX.style.color = '#000000';
    }

    let queryY = document.body.querySelector('#yInput'); 
    if (isNaN(y)) {
      queryY.style.color = '#D70040'; ready = false;
      ready = false;
    } else {
      queryY.style.color = '#000000';
    }
    if (ready) {
      document.body.querySelectorAll('button').forEach((element) =>{
        element.style.disabled = 'true'
      })
      this.state.isMergeEnabled ? 
        StartMergeProcess(this.state.dir1, this.state.dir2, this.state.dir3) : 
        StartCropProcess(this.state.dir1, this.state.dir3, x, y)
    }

    //StartMergeProcess('test', 'test', 'test');
  };

  render() {
    return (
      <div id='App'>
        <div className='selectDataset'>
          <div className='operation'><h2>Operation</h2></div>
          <p></p>
          <div>
            <input
              type='radio'
              id='crop'
              name='rad'
              checked={!this.state.isMergeEnabled}
              onChange={this.ChangeMode}
            />
            <label htmlFor='crop' className='radios'>
              Crop
            </label>
            <input
              type='radio'
              id='merge'
              name='rad'
              checked={this.state.isMergeEnabled}
              onChange={this.ChangeMode}
            />
            <label htmlFor='merge' className='radios'>
              Merge
            </label>
          </div>
          <p></p>
          <div
            className={`selectDirs ${
              this.state.isMergeEnabled ? '' : 'cropEnabled'
            }`}
          >
            <div className='selectDirField'>
              <p id='selectdir1' className='selectdir1'>
                {this.state.dir1}
              </p>
              <button className='btn' onClick={this.SetDirectory1}>
                Select
              </button>
            </div>
            <div
              className={`${this.state.isMergeEnabled ? '' : 'hide'}`}
              style={{color: '#A9A9A9'}}
            >
              +
            </div>
            <div
              className={`${this.state.isMergeEnabled ? 'hide' : ''}`}
              style={{color: '#A9A9A9'}}
            >
              ↓
            </div>
            <div
              className={`selectDirField ${
                this.state.isMergeEnabled ? '' : 'hide'
              }`}
            >
              <p id='selectdir2'>{this.state.dir2}</p>
              <button className='btn' onClick={this.SetDirectory2}>
                Select
              </button>
            </div>
            <div
              className={`${this.state.isMergeEnabled ? '' : 'hide'}`}
              style={{color: '#A9A9A9'}}
            >
              =
            </div>
            <div className='selectDirField'>
              <p id='selectdir3' className='selectdir3'>
                {this.state.dir3}
              </p>
              <button className='btn' onClick={this.SetDirectory3}>
                Select
              </button>
            </div>
            <div className={`selectDirField sizes ${
                !this.state.isMergeEnabled ? '' : 'hide'
              }`}>
            <div>Size (px)</div>
            <div>x: <input id='xInput' maxLength='5' defaultValue={64}></input></div>
            <div>y: <input id='yInput' maxLength='5' defaultValue={64}></input></div>
            </div>
          </div>
          <div>
            <button className='start' onClick={this.Progress}>Start</button>
          </div>
        </div>
      </div>
    );
  }
}
