"use client";

import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { FileUpload } from "@/components/ui/file-upload";
import { 
  getAllDocuments, 
  createDocument, 
  updateDocument, 
  deleteDocument, 
  getDocumentStats,
  downloadDocument
} from "@/lib/api";
import { PolicyFile, CreateDocumentRequest, UpdateDocumentRequest, DocumentStats } from "@/lib/types";
import { 
  Plus, 
  Search, 
  Edit, 
  Trash2, 
  FileText, 
  Users, 
  Filter,
  ArrowLeft,
  BarChart3,
  Clock,
  Download
} from "lucide-react";
import Link from "next/link";
import { DocumentManagementRoute } from "@/components/auth/ProtectedRoute";
import { useAuth } from "@/lib/auth";
import UserManagement from "@/components/admin/UserManagement";
import AuditLogs from "@/components/admin/AuditLogs";

export default function Dashboard() {
  return (
    <DocumentManagementRoute>
      <DashboardContent />
    </DocumentManagementRoute>
  );
}

function DashboardContent() {
  const [activeTab, setActiveTab] = useState<"documents" | "users" | "audit">("documents");
  const [searchQuery, setSearchQuery] = useState("");
  const [filterType, setFilterType] = useState<"all" | "policy" | "onboarding">("all");
  const [filterCategory, setFilterCategory] = useState("");
  const [selectedDocument, setSelectedDocument] = useState<PolicyFile | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  const { user } = useAuth();
  const queryClient = useQueryClient();

  // Fetch documents
  const { data: documents = [], isLoading } = useQuery({
    queryKey: ['documents', filterType, filterCategory],
    queryFn: () => getAllDocuments({
      type: filterType === "all" ? undefined : filterType,
      category: filterCategory || undefined,
      active: true
    }),
  });

  // Fetch stats
  const { data: stats } = useQuery({
    queryKey: ['document-stats'],
    queryFn: getDocumentStats,
  });

  // Create document mutation
  const createMutation = useMutation({
    mutationFn: createDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      queryClient.invalidateQueries({ queryKey: ['document-stats'] });
      setIsCreateDialogOpen(false);
    },
  });

  // Update document mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, updates }: { id: number; updates: UpdateDocumentRequest }) => updateDocument(id, updates),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      queryClient.invalidateQueries({ queryKey: ['document-stats'] });
      setIsEditDialogOpen(false);
      setSelectedDocument(null);
    },
  });

  // Delete document mutation
  const deleteMutation = useMutation({
    mutationFn: deleteDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      queryClient.invalidateQueries({ queryKey: ['document-stats'] });
    },
  });

  // Filter documents based on search query
  const filteredDocuments = documents.filter(doc =>
    doc.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    doc.content.toLowerCase().includes(searchQuery.toLowerCase()) ||
    doc.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
    doc.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  // Get unique categories for filtering
  const categories = Array.from(new Set(documents.map(doc => doc.category)));

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Link href="/">
              <Button variant="outline" size="sm">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back to Chat
              </Button>
            </Link>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
              <p className="text-gray-600">
                Manage documents, users, and system settings
                {user && (
                  <span className="ml-2 text-sm">
                    • Logged in as <strong>{user.first_name} {user.last_name}</strong> ({user.role})
                  </span>
                )}
              </p>
            </div>
          </div>
          {activeTab === "documents" && (
            <Button 
              onClick={() => setIsCreateDialogOpen(true)}
              className="bg-blue-600 hover:bg-blue-700"
            >
              <Plus className="w-4 h-4 mr-2" />
              Add Document
            </Button>
          )}
        </div>

        {/* Tab Navigation */}
        <div className="flex space-x-1 mb-8">
          <Button
            variant={activeTab === "documents" ? "default" : "outline"}
            onClick={() => setActiveTab("documents")}
            className="flex items-center gap-2"
          >
            <FileText className="w-4 h-4" />
            Documents
          </Button>
          <Button
            variant={activeTab === "users" ? "default" : "outline"}
            onClick={() => setActiveTab("users")}
            className="flex items-center gap-2"
          >
            <Users className="w-4 h-4" />
            Users
          </Button>
          <Button
            variant={activeTab === "audit" ? "default" : "outline"}
            onClick={() => setActiveTab("audit")}
            className="flex items-center gap-2"
          >
            <Clock className="w-4 h-4" />
            Audit Logs
          </Button>
        </div>

        {/* Tab Content */}
        {activeTab === "documents" && (
          <DocumentsTab
            documents={filteredDocuments}
            stats={stats}
            isLoading={isLoading}
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            filterType={filterType}
            setFilterType={setFilterType}
            filterCategory={filterCategory}
            setFilterCategory={setFilterCategory}
            categories={categories}
            onEditDocument={(doc) => {
              setSelectedDocument(doc);
              setIsEditDialogOpen(true);
            }}
            onDeleteDocument={(id) => deleteMutation.mutate(id)}
            onDownloadDocument={async (doc) => {
              try {
                await downloadDocument(doc.id);
              } catch (error) {
                console.error('Download error:', error);
              }
            }}
          />
        )}

        {activeTab === "users" && (
          <UserManagement />
        )}

        {activeTab === "audit" && (
          <AuditLogs />
        )}

        {/* Create Document Dialog */}
        <CreateDocumentDialog
          open={isCreateDialogOpen}
          onOpenChange={setIsCreateDialogOpen}
          onSubmit={(data) => createMutation.mutate(data)}
          isLoading={createMutation.isPending}
        />

        {/* Edit Document Dialog */}
        {selectedDocument && (
          <EditDocumentDialog
            open={isEditDialogOpen}
            onOpenChange={setIsEditDialogOpen}
            document={selectedDocument}
            onSubmit={(updates) => updateMutation.mutate({ id: selectedDocument.id, updates })}
            isLoading={updateMutation.isPending}
          />
        )}
      </div>
    </div>
  );
}

// Documents Tab Component
function DocumentsTab({
  documents,
  stats,
  isLoading,
  searchQuery,
  setSearchQuery,
  filterType,
  setFilterType,
  filterCategory,
  setFilterCategory,
  categories,
  onEditDocument,
  onDeleteDocument,
  onDownloadDocument,
}: {
  documents: PolicyFile[];
  stats?: DocumentStats;
  isLoading: boolean;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  filterType: "all" | "policy" | "onboarding";
  setFilterType: (type: "all" | "policy" | "onboarding") => void;
  filterCategory: string;
  setFilterCategory: (category: string) => void;
  categories: string[];
  onEditDocument: (doc: PolicyFile) => void;
  onDeleteDocument: (id: number) => void;
  onDownloadDocument: (doc: PolicyFile) => void;
}) {
  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="animate-pulse">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-24 bg-gray-200 rounded"></div>
            ))}
          </div>
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-20 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <>
      {/* Statistics Cards */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Documents</CardTitle>
              <FileText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.total}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Policies</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.policies}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Onboarding Docs</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.onboardingDocs}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Categories</CardTitle>
              <Filter className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{Object.keys(stats.categories).length}</div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4 mb-6">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
          <Input
            placeholder="Search documents..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9"
          />
        </div>
        
        <div className="flex gap-2">
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value as "all" | "policy" | "onboarding")}
            className="px-3 py-2 border border-gray-300 rounded-md text-sm"
          >
            <option value="all">All Types</option>
            <option value="policy">Policies</option>
            <option value="onboarding">Onboarding</option>
          </select>
          
          <select
            value={filterCategory}
            onChange={(e) => setFilterCategory(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md text-sm"
          >
            <option value="">All Categories</option>
            {categories.map((category) => (
              <option key={category} value={category}>
                {category}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Documents Grid */}
      <div className="grid gap-6">
        {documents.length > 0 ? (
          documents.map((doc) => (
            <Card key={doc.id} className="hover:shadow-md transition-shadow">
              <CardContent className="p-6">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <h3 className="text-lg font-semibold">{doc.name}</h3>
                      <Badge variant={doc.document_type === 'policy' ? 'default' : 'secondary'}>
                        {doc.document_type}
                      </Badge>
                      {!doc.is_active && (
                        <Badge variant="destructive">Inactive</Badge>
                      )}
                    </div>
                    <p className="text-gray-600 text-sm mb-2">{doc.description}</p>
                                         <div className="flex items-center gap-4 text-xs text-gray-500">
                       <span>Category: {doc.category}</span>
                       <span>Updated: {new Date(doc.last_updated).toLocaleDateString()}</span>
                       <span>By: {doc.created_by}</span>
                     </div>
                    {doc.tags && doc.tags.length > 0 && (
                      <div className="flex flex-wrap gap-1 mt-2">
                        {doc.tags.map((tag, index) => (
                          <Badge key={index} variant="outline" className="text-xs">
                            {tag}
                          </Badge>
                        ))}
                      </div>
                    )}
                  </div>
                  
                  <div className="flex gap-2 ml-4">
                    {doc.file_path && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => onDownloadDocument(doc)}
                        title="Download original file"
                      >
                        <Download className="w-4 h-4" />
                      </Button>
                    )}
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => onEditDocument(doc)}
                    >
                      <Edit className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        if (confirm(`Are you sure you want to delete "${doc.name}"?`)) {
                          onDeleteDocument(doc.id);
                        }
                      }}
                      className="text-red-600 hover:text-red-700"
                    >
                      <Trash2 className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
                <div className="text-sm text-gray-700 bg-gray-50 p-3 rounded-md">
                  {doc.content.length > 200 
                    ? `${doc.content.substring(0, 200)}...` 
                    : doc.content
                  }
                </div>
              </CardContent>
            </Card>
          ))
        ) : (
          <div className="text-center py-12">
            <FileText className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">No documents found</h3>
            <p className="text-gray-500 mb-4">
              {searchQuery 
                ? "No documents match your search criteria" 
                : "Get started by creating your first document"}
            </p>
          </div>
        )}
      </div>
    </>
  );
}

// Create Document Dialog Component
function CreateDocumentDialog({
  open,
  onOpenChange,
  onSubmit,
  isLoading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: CreateDocumentRequest) => void;
  isLoading: boolean;
}) {
  const [formData, setFormData] = useState<CreateDocumentRequest>({
    name: '',
    content: '',
    description: '',
    category: '',
    document_type: 'policy',
    tags: [],
    created_by: '',
  });

  const [currentTag, setCurrentTag] = useState('');
  const [inputMethod, setInputMethod] = useState<'manual' | 'upload'>('manual');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const addTag = () => {
    if (currentTag.trim() && !(formData.tags || []).includes(currentTag.trim())) {
      setFormData(prev => ({
        ...prev,
        tags: [...(prev.tags || []), currentTag.trim()]
      }));
      setCurrentTag('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: (prev.tags || []).filter(tag => tag !== tagToRemove)
    }));
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addTag();
    }
  };

  const handleFileUpload = (data: {
    fileName: string;
    fileType: string;
    fileSize: number;
    filePath?: string;
    extractedText: string;
  }) => {
    // Clean the filename by removing extension and timestamp pattern
    let cleanName = data.fileName.replace(/\.[^/.]+$/, ""); // Remove extension
    
    // Remove timestamp pattern: _YYYYMMDD_HHMMSS
    const timestampPattern = /_\d{8}_\d{6}$/;
    if (timestampPattern.test(cleanName)) {
      cleanName = cleanName.replace(timestampPattern, "");
    }
    
    // Auto-fill form with extracted data
    setFormData(prev => ({
      ...prev,
      name: prev.name || cleanName,
      content: data.extractedText,
      description: prev.description || `Document uploaded from ${data.fileName}`,
      file_path: data.filePath,
    }));
  };

  const handleFileUploadError = (error: string) => {
    console.error('File upload error:', error);
    // You could add a toast notification here
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Create New Document</DialogTitle>
        </DialogHeader>
        
        {/* Input Method Toggle */}
        <div className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
          <span className="text-sm font-medium">Input Method:</span>
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => setInputMethod('manual')}
              className={`px-3 py-1 text-sm rounded ${
                inputMethod === 'manual'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              Manual Entry
            </button>
            <button
              type="button"
              onClick={() => setInputMethod('upload')}
              className={`px-3 py-1 text-sm rounded ${
                inputMethod === 'upload'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              File Upload
            </button>
          </div>
        </div>

        {/* File Upload Section */}
        {inputMethod === 'upload' && (
          <div className="space-y-4">
            <FileUpload
              onFileUploaded={handleFileUpload}
              onError={handleFileUploadError}
              maxSize={10 * 1024 * 1024} // 10MB
            />
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Document Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                placeholder="e.g., Password Policy"
                required
              />
            </div>
            <div>
              <label className="text-sm font-medium">Document Type</label>
              <select
                value={formData.document_type}
                onChange={(e) => setFormData(prev => ({ ...prev, document_type: e.target.value as 'policy' | 'onboarding' }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              >
                <option value="policy">Policy</option>
                <option value="onboarding">Onboarding</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Category</label>
              <Input
                value={formData.category}
                onChange={(e) => setFormData(prev => ({ ...prev, category: e.target.value }))}
                placeholder="e.g., Security, HR, IT"
                required
              />
            </div>
            <div>
              <label className="text-sm font-medium">Created By</label>
              <Input
                value={formData.created_by}
                onChange={(e) => setFormData(prev => ({ ...prev, created_by: e.target.value }))}
                placeholder="Author name"
                required
              />
            </div>
          </div>

          <div>
            <label className="text-sm font-medium">Description</label>
            <Input
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              placeholder="Brief description of the document"
              required
            />
          </div>

          <div>
            <label className="text-sm font-medium">
              Content {inputMethod === 'upload' && '(Auto-filled from uploaded file)'}
            </label>
            <textarea
              value={formData.content}
              onChange={(e) => setFormData(prev => ({ ...prev, content: e.target.value }))}
              placeholder={inputMethod === 'upload' ? "Content will be filled automatically when you upload a file..." : "Document content..."}
              className="w-full px-3 py-2 border border-gray-300 rounded-md min-h-[120px]"
              readOnly={inputMethod === 'upload'}
              required
            />
          </div>

          <div>
            <label className="text-sm font-medium">Tags</label>
            <div className="flex gap-2 mb-2">
              <Input
                value={currentTag}
                onChange={(e) => setCurrentTag(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Add a tag"
              />
              <Button type="button" onClick={addTag} variant="outline">
                Add
              </Button>
            </div>
            <div className="flex flex-wrap gap-1">
              {(formData.tags || []).map((tag, index) => (
                <Badge key={index} variant="secondary" className="cursor-pointer" onClick={() => removeTag(tag)}>
                  {tag} ×
                </Badge>
              ))}
            </div>
          </div>



          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? 'Creating...' : 'Create Document'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

// Edit Document Dialog Component
function EditDocumentDialog({
  open,
  onOpenChange,
  document,
  onSubmit,
  isLoading,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  document: PolicyFile;
  onSubmit: (updates: UpdateDocumentRequest) => void;
  isLoading: boolean;
}) {
  const [formData, setFormData] = useState<UpdateDocumentRequest>({
    name: document.name,
    content: document.content,
    description: document.description,
    category: document.category,
    document_type: document.document_type,
    tags: document.tags || [],
    is_active: document.is_active,
  });

  const [currentTag, setCurrentTag] = useState('');

  useEffect(() => {
    setFormData({
      name: document.name,
      content: document.content,
      description: document.description,
      category: document.category,
      document_type: document.document_type,
      tags: document.tags || [],
      is_active: document.is_active,
    });
  }, [document]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const addTag = () => {
    if (currentTag.trim() && !(formData.tags || []).includes(currentTag.trim())) {
      setFormData((prev: UpdateDocumentRequest) => ({
        ...prev,
        tags: [...(prev.tags || []), currentTag.trim()]
      }));
      setCurrentTag('');
    }
  };

  const removeTag = (tagToRemove: string) => {
    setFormData((prev: UpdateDocumentRequest) => ({
      ...prev,
      tags: (prev.tags || []).filter(tag => tag !== tagToRemove)
    }));
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addTag();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Document</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Document Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData((prev: UpdateDocumentRequest) => ({ ...prev, name: e.target.value }))}
                required
              />
            </div>
            <div>
              <label className="text-sm font-medium">Document Type</label>
              <select
                value={formData.document_type}
                onChange={(e) => setFormData((prev: UpdateDocumentRequest) => ({ ...prev, document_type: e.target.value as 'policy' | 'onboarding' }))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md"
                required
              >
                <option value="policy">Policy</option>
                <option value="onboarding">Onboarding</option>
              </select>
            </div>
          </div>

          <div>
            <label className="text-sm font-medium">Category</label>
            <Input
              value={formData.category}
              onChange={(e) => setFormData((prev: UpdateDocumentRequest) => ({ ...prev, category: e.target.value }))}
              required
            />
          </div>

          <div>
            <label className="text-sm font-medium">Description</label>
            <Input
              value={formData.description}
              onChange={(e) => setFormData((prev: UpdateDocumentRequest) => ({ ...prev, description: e.target.value }))}
              required
            />
          </div>

          <div>
            <label className="text-sm font-medium">Content</label>
            <textarea
              value={formData.content}
              onChange={(e) => setFormData((prev: UpdateDocumentRequest) => ({ ...prev, content: e.target.value }))}
              className="w-full px-3 py-2 border border-gray-300 rounded-md min-h-[120px]"
              required
            />
          </div>

          <div>
            <label className="text-sm font-medium">Tags</label>
            <div className="flex gap-2 mb-2">
              <Input
                value={currentTag}
                onChange={(e) => setCurrentTag(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Add a tag"
              />
              <Button type="button" onClick={addTag} variant="outline">
                Add
              </Button>
            </div>
                         <div className="flex flex-wrap gap-1">
               {(formData.tags || []).map((tag, index) => (
                 <Badge key={index} variant="secondary" className="cursor-pointer" onClick={() => removeTag(tag)}>
                   {tag} ×
                 </Badge>
               ))}
             </div>
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="is_active"
              checked={formData.is_active}
              onChange={(e) => setFormData((prev: UpdateDocumentRequest) => ({ ...prev, is_active: e.target.checked }))}
            />
            <label htmlFor="is_active" className="text-sm font-medium">
              Active Document
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? 'Updating...' : 'Update Document'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
} 