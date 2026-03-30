import { type NavigateFunction } from 'react-router-dom';


export function BackButton({ retTo, navi }: { retTo: string; navi:NavigateFunction}) {
    
    const goBack = () => {
        navi(retTo);
    }    

    return (
        <button className="monitor-button" onClick={() => goBack()}>
          {"<<"}
        </button>
    );
};