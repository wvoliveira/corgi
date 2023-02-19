import React from "react";

const TagInput = ({ tagList, addTag, removeTag }) => {
  const [tag, setTag] = React.useState("");

  const changeTagInput = (e) => setTag(e.target.value);

  const handleTagInputKeyDown = (e) => {
    switch (e.keyCode) {
      case 13: // Enter
      case 9: // Tab
      case 188: // Comma
        if (e.keyCode !== 9) e.preventDefault();
        handleAddTag();
        break;
      default:
        break;
    }
  };

  const handleAddTag = () => {
    if (!!tag) {
      addTag(tag);
      setTag("");
    }
  };

  const handleRemoveTag = (tag) => {
    removeTag(tag);
  };

  return (
    <>
      <input
        type="text"
        placeholder="Enter tags"
        value={tag}
        onChange={changeTagInput}
        onBlur={handleAddTag}
        onKeyDown={handleTagInputKeyDown}
      />

      <div>
        {tagList.map((tag, index) => (
          <span key={index}>
            <i
              onClick={() => handleRemoveTag(tag)}
            />
            {tag}
          </span>
        ))}
      </div>
    </>
  );
};

export default TagInput;
