export class Session {
  private readonly baseURL: string = '';
  private readonly commands: string[] = [];
  private _onmessage: (string) => void = (_) => {};
  private _onerror: (Error) => void = (err) => console.log(err);
  private socket: WebSocket;

  constructor() {
    if (process.env.NODE_ENV == 'development') {
      this.baseURL = 'http://localhost:3030';
    }
  }

  metaUrl(): string {
    return `${this.baseURL}/meta`;
  }

  set onmessage(value: (str: string) => void) {
    this._onmessage = value;
  }

  set onerror(value: (err: Error) => void) {
    this._onerror = value;
  }

  connect() {
    let wsUrl = '';
    if (this.baseURL == '') {
      wsUrl = 'ws://' + window.location.host + '/ws';
      console.log(wsUrl);
    } else {
      wsUrl = 'ws://localhost:3030/ws';
    }
    this.socket = new WebSocket(wsUrl);
    this.socket.onmessage = (msg) => this._onmessage(msg.data);
    this.socket.onerror = (err) => this._onerror(err);
    this.socket.onopen = (ev) => {
      this.flush();
    };
  }

  flush() {
    while (this.commands.length > 0) {
      const cmd = this.commands.shift();
      this.socket.send(cmd);
    }
  }

  command(cmd: string) {
    if (this.socket.readyState == 1) {
      this.flush();
      this.socket.send(cmd);
    } else {
      this.commands.push(cmd);
    }
  }
}
