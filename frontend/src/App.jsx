import * as React from 'react';
import './App.css';
import { SetDirectoryDialog } from '../wailsjs/go/main/App';
import { StartMergeProcess } from '../wailsjs/go/main/App';
import { StartCropProcess } from '../wailsjs/go/main/App';
import { StartProcessBrightness } from '../wailsjs/go/main/App' 

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      operation: 0,
      checked: [1, 0, 0],
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
    let flag = 0;
    let checked = [0, 0, 0]
    if (document.body.querySelector('#crop').checked) {
      checked = [1, 0, 0]
      flag = 0;
    }
    if (document.body.querySelector('#merge').checked) {
      checked = [0, 1, 0]
      flag = 1;
    }
    if (document.body.querySelector('#brightness').checked) {
      checked = [0, 0, 1]
      flag = 2;
    }
    let dirText = this.state.dir1;
    if (!this.state.dir1Selected) {
      dirText = flag ? 'Files directory' : 'First files directory';
    }
    document.body.querySelector('#selectdir1').style.color = '#FFFFFF'
    document.body.querySelector('#selectdir2').style.color = '#FFFFFF'
    document.body.querySelector('#selectdir3').style.color = '#FFFFFF'
    this.setState({ checked: checked, operation: flag, dir1: dirText });
  };

  Progress = () => {
    console.log(this.state.operation)
    let ready = true;
    if (!this.state.dir1Selected) {document.body.querySelector('#selectdir1').style.color = '#D70040'; ready = false;}
    if (!this.state.dir2Selected && this.state.operation == 1) {document.body.querySelector('#selectdir2').style.color = '#D70040'; ready = false;}
    if (!this.state.dir3Selected) {document.body.querySelector('#selectdir3').style.color = '#D70040'; ready = false;}
    
    let x = document.getElementById('xInput').value;
    let y = document.getElementById('yInput').value;
    let br = document.getElementById('brInput').value;
    x = parseInt(x); 
    y = parseInt(y);
    br = parseInt(br);

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

    let queryBr = document.body.querySelector('#brInput'); 
    if (isNaN(br)) {
      queryBr.style.color = '#D70040'; ready = false;
      ready = false;
    } else {
      queryBr.style.color = '#000000';
    }
    
    if (ready) {
      document.body.querySelectorAll('button').forEach((element) =>{
        element.style.backgroundColor = '#808080'
        element.disabled = true;
      })

      switch (this.state.operation) {
        case 0: {StartCropProcess(this.state.dir1, this.state.dir3, x, y); break;}
        case 1: {StartMergeProcess(this.state.dir1, this.state.dir2, this.state.dir3);  break;}
        case 2: {StartProcessBrightness(this.state.dir1, this.state.dir3, br); break;}
      }


      document.body.querySelectorAll('button').forEach((element) =>{
        element.style.backgroundColor = '#2ea44f'
        element.disabled = false;
      })
    }
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
              checked={this.state.checked[0]}
              onChange={this.ChangeMode}
            />
            <label htmlFor='crop' className='radios'>
              Crop
            </label>
            <input
              type='radio'
              id='merge'
              name='rad'
              checked={this.state.checked[1]}
              onChange={this.ChangeMode}
            />
            <label htmlFor='merge' className='radios'>
              Merge
            </label>
            <input
              type='radio'
              id='brightness'
              name='rad'
              checked={this.state.checked[2]}
              onChange={this.ChangeMode}
            />
            <label htmlFor='brightness' className='radios'>
              Brightness+
            </label>
          </div>
          <p></p>
          <div
            className={`selectDirs ${
              this.state.operation === 0 ? 'cropEnabled' : ''
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
              className={`${this.state.operation === 1 ? '' : 'hide'}`}
              style={{color: '#A9A9A9'}}
            >
              +
            </div>
            <div
              className={`${this.state.operation !== 1 ? '' : 'hide'}`}
              style={{color: '#A9A9A9'}}
            >
              ↓
            </div>
            <div
              className={`selectDirField ${
                this.state.operation === 1 ? '' : 'hide'
              }`}
            >
              <p id='selectdir2'>{this.state.dir2}</p>
              <button className='btn' onClick={this.SetDirectory2}>
                Select
              </button>
            </div>
            <div
              className={`${this.state.operation === 1 ? '' : 'hide'}`}
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
                this.state.operation === 0 ? '' : 'hide'
              }`}>
            <div>Size (px)</div>
            <div>x: <input id='xInput' maxLength='5' defaultValue={64}></input></div>
            <div>y: <input id='yInput' maxLength='5' defaultValue={64}></input></div>
            </div>

            <div className={`selectDirField sizes ${
                this.state.operation === 2 ? '' : 'hide'
              }`}>
            <div>Increaase brightness by</div>
            <div><input id='brInput' maxLength='5' defaultValue={8}></input></div>
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
