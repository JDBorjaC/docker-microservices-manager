import '../styles/monitor.css';
import monitor_url from '../assets/msmmonitor.svg';

interface MonitorBackdropProps {
    children: React.ReactNode;
}

export default function MonitorBackdrop({ children }: MonitorBackdropProps) {
  return (
    <div className="wrapper">
      <img src={monitor_url} alt="monitor" className="monitor-img" />
      <div className="monitor-hole">
        {children}
      </div>
    </div>
  );
}