import { useState } from "react";


export function ExpandableMessage({ title, children }: { title: string; children: React.ReactNode }) {
  const [open, setOpen] = useState(false);

  return (
    <div>
      <button className="expand-button" onClick={() => setOpen(prev => !prev)}>
        {open ? '▼' : '▶'} {title}
      </button>
      {open && <div className="expand-text">{children}</div>}
    </div>
  );
};