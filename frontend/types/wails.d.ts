export interface UserConfig{
  username: string;
}

export interface FileInfo {
  name: string;
  path: string;
  size: number;
}

export interface MainInterface {
  SelectFile(): Promise<FileInfo>;
  SendFile(file: FileInfo): Promise<void>;
  ReceiveFile(): Promise<void>;
}

declare global {
  interface Window {
    go: {
      main: {
        App: MainInterface;
      };
    };
  }
}

export {};
