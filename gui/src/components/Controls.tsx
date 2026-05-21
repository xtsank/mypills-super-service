import { ChangeEvent } from "react";

type BaseFieldProps = {
  label: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  type?: "text" | "password";
  readOnly?: boolean;
};

export function TextInput({ label, value, onChange, placeholder, type = "text", readOnly }: BaseFieldProps) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <input
        className="input"
        type={type}
        value={value}
        onChange={(event: ChangeEvent<HTMLInputElement>) => onChange(event.target.value)}
        placeholder={placeholder}
        readOnly={readOnly}
      />
    </label>
  );
}

export function NumberInput({ label, value, onChange, placeholder }: BaseFieldProps) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <input
        className="input"
        type="number"
        value={value}
        onChange={(event: ChangeEvent<HTMLInputElement>) => onChange(event.target.value)}
        placeholder={placeholder}
      />
    </label>
  );
}

export function Checkbox({ label, checked, onChange }: { label: string; checked: boolean; onChange: (value: boolean) => void }) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <input
        type="checkbox"
        checked={checked}
        onChange={(event) => onChange(event.target.checked)}
      />
    </label>
  );
}

export function TextArea({ label, value, onChange, placeholder }: BaseFieldProps) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <textarea
        className="textarea"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
      />
    </label>
  );
}

export function SelectInput({
  label,
  value,
  options,
  onChange,
  placeholder
}: {
  label: string;
  value: string;
  options: Array<{ value: string; label: string }>;
  onChange: (value: string) => void;
  placeholder?: string;
}) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <select
        className="select"
        value={value}
        onChange={(event) => onChange(event.target.value)}
      >
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}

export function MultiSelectInput({
  label,
  values,
  options,
  onChange,
  placeholder
}: {
  label: string;
  values: string[];
  options: Array<{ value: string; label: string }>;
  onChange: (value: string[]) => void;
  placeholder?: string;
}) {
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <select
        className="select"
        multiple
        value={values}
        onChange={(event) =>
          onChange(Array.from(event.currentTarget.selectedOptions).map((option) => option.value))
        }
      >
        {placeholder && options.length === 0 && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}

export function Button({
  children,
  onClick,
  variant = "primary",
  type = "button"
}: {
  children: string;
  onClick?: () => void;
  variant?: "primary" | "secondary";
  type?: "button" | "submit";
}) {
  const className = variant === "secondary" ? "button button--secondary" : "button";
  return (
    <button type={type} className={className} onClick={onClick}>
      {children}
    </button>
  );
}

export function TagInput({
  label,
  value,
  onChange,
  placeholder
}: {
  label: string;
  value: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
}) {
  const text = value.join(", ");
  return (
    <label className="field">
      <span className="field__label">{label}</span>
      <input
        className="input"
        value={text}
        onChange={(event) =>
          onChange(
            event.target.value
              .split(",")
              .map((item) => item.trim())
              .filter((item) => item.length > 0)
          )
        }
        placeholder={placeholder}
      />
    </label>
  );
}
