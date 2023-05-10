from dataclasses import dataclass
from typing import Optional


@dataclass
class FileUploadInput(object):
    """Class represents the file_upload.input"""
    space_id: str
    name: Optional[str]
    file_id: str
    parent_folder_id: str
    contents: dict
    if_match: Optional[str]
    if_none_match: Optional[str]


@dataclass
class FileUploadResult(object):
    """Class represents the file_upload.result"""
    version: Optional[int]


@dataclass
class FileUpload(object):
    """Class represents the file_upload collection"""
    id: str
    status: str
    input: FileUploadInput
    result: Optional[FileUploadResult]
    created_on: str
    created_by: dict
    modified_on: str
    error_reason: Optional[str]
