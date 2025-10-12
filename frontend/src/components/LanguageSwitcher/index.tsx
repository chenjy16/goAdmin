import React from 'react';
import { Select } from 'antd';
import { GlobalOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

const { Option } = Select;

interface LanguageSwitcherProps {
  className?: string;
}

const LanguageSwitcher: React.FC<LanguageSwitcherProps> = ({ className }) => {
  const { i18n, t } = useTranslation();

  const handleLanguageChange = (language: string) => {
    i18n.changeLanguage(language);
  };

  return (
    <Select
      value={i18n.language}
      onChange={handleLanguageChange}
      className={className}
      style={{ width: 120 }}
      suffixIcon={<GlobalOutlined />}
    >
      <Option value="zh">{t('settings.chinese')}</Option>
      <Option value="en">{t('settings.english')}</Option>
    </Select>
  );
};

export default LanguageSwitcher;