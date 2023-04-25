export interface message {
    id : number;
    id_histori : number;
    Jenis : "input" | "output";
    Isi : string;
}

export interface chatProps {
    className: string;
    messages: message[];
}

export interface messageProps {
    message: message;
}

export interface sidebarProps {
    className: string,
    histories: history[]
  }

export interface history{
    id: number,
    nama: string,
}

export interface buttonProps{
    history: history,
    clicked: number,
    handleClick: (i:number) => void
    handleDelete: (i:number) => void
}

export interface sendMessageProps{
    inputValue: string,
    setInputValue: React.Dispatch<React.SetStateAction<string>>,
    handleInput: (e: React.FormEvent<HTMLFormElement>) => void
}